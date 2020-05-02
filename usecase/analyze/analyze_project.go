package analyze

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase/analyze/step"
	log "github.com/sirupsen/logrus"
)

var (
	ErrUnableToBuildASTs        = errors.New("unable to create ASTs from input")
	ErrUnableToMineASTs         = errors.New("unable to apply one or more miners to the ASTs")
	ErrUnableToCreateProcessors = errors.New("unable to create splitting or expansion algorithms")
	ErrUnableToSaveIdentifiers  = errors.New("unable to save extracted and processed indentifiers")
)

type AnalyzeProjectUsecase interface {
	Analyze(ctx context.Context, project entity.Project, config *entity.AnalysisConfig) (entity.AnalysisResults, error)
}

func NewAnalyzeProjectUsecase(scr repository.SourceCodeRepository, ir repository.IdentifierRepository) AnalyzeProjectUsecase {
	return &analyzeProjectUsecase{
		sourceCodeRepository: scr,
		identifierRepository: ir,
	}
}

type analyzeProjectUsecase struct {
	sourceCodeRepository repository.SourceCodeRepository
	identifierRepository repository.IdentifierRepository
}

func (uc analyzeProjectUsecase) Analyze(ctx context.Context, project entity.Project, config *entity.AnalysisConfig) (entity.AnalysisResults, error) {
	analysisResults := entity.AnalysisResults{
		ID:                project.ID,
		DateCreated:       time.Now(),
		ProjectName:       project.Metadata.Fullname,
		ProjectURL:        project.URL,
		PipelineMiners:    make([]string, 0),
		PipelineSplitters: make([]string, 0),
		PipelineExpanders: make([]string, 0),
	}
	// read and parse files
	filesc := step.Read(ctx, uc.sourceCodeRepository, project.SourceCode.Location, project.SourceCode.Files)
	parsed := step.Parse(filesc)
	files := step.Merge(parsed)

	valid := make([]entity.File, 0)
	fileErrorSamples := make([]string, 0)
	for _, file := range files {
		if file.Error != nil {
			if len(fileErrorSamples) < 10 {
				fileErrorSamples = append(fileErrorSamples, file.Error.Error())
			}

			log.WithError(file.Error).Warn(fmt.Sprintf("unable to read or parse file %s at %s", file.Name, project.SourceCode.Location))
			continue
		}
		valid = append(valid, file)
	}
	analysisResults.FilesTotal = len(files)
	analysisResults.FilesValid = len(valid)
	analysisResults.FilesError = len(files) - len(valid)
	analysisResults.FilesErrorSamples = fileErrorSamples

	// if every file can't be parsed, then fail
	if len(valid) == 0 {
		log.Error(fmt.Sprintf("unable to read or parse any file on %s for project %s",
			project.SourceCode.Location, project.URL))
		return entity.AnalysisResults{}, ErrUnableToBuildASTs
	}

	// apply the pre-process step (mine them)
	miners := buildMiners(config)
	for _, miner := range miners {
		analysisResults.PipelineMiners = append(analysisResults.PipelineMiners, miner.Name())
	}

	miningResults := step.Mine(valid, miners...)

	// make the splitters from input and mining results
	splitters := buildSplittersFromInputAndMiningResults(config, miningResults)
	if len(splitters) == 0 {
		log.WithField("desired", config.Splitters).Error("unable to create any splitter")
		return entity.AnalysisResults{}, ErrUnableToCreateProcessors
	}
	for _, splitter := range splitters {
		analysisResults.PipelineSplitters = append(analysisResults.PipelineSplitters, splitter.Name())
	}

	// make the expanders from input and mining results
	expanders := buildExpandersFromInputAndMiningResults(config, miningResults)
	if len(expanders) == 0 {
		log.WithField("desired", config.Expanders).Error("unable to create any expander")
		return entity.AnalysisResults{}, ErrUnableToCreateProcessors
	}
	for _, expander := range expanders {
		analysisResults.PipelineExpanders = append(analysisResults.PipelineExpanders, expander.Name())
	}

	// analyze each identifier
	identc := step.Extract(valid, config.ExtractorFactory)
	splittedc := step.Split(identc, splitters...)
	expandedc := step.Expand(splittedc, expanders...)

	identErrorSamples := make([]string, 0)
	for ident := range expandedc {
		analysisResults.IdentifiersTotal++
		if ident.Error != nil {
			if len(identErrorSamples) < 10 {
				identErrorSamples = append(identErrorSamples, ident.Error.Error())
			}

			log.WithFields(log.Fields{
				"project":    project.Metadata.Fullname,
				"filename":   ident.File,
				"identifier": ident.Name,
				"position":   ident.Position,
			}).Warn("an error occurred during the splitting the expansion for the identifier")
			analysisResults.IdentifiersError++
		}

		err := uc.identifierRepository.Add(ctx, project, ident)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to save identifier %s, on file %s for project %s",
				ident.Name, ident.File, project.URL))
			return entity.AnalysisResults{}, ErrUnableToSaveIdentifiers
		}
	}
	analysisResults.IdentifiersValid = analysisResults.IdentifiersTotal - analysisResults.IdentifiersError
	analysisResults.IdentifiersErrorSamples = identErrorSamples

	return analysisResults, nil
}

func buildMiners(config *entity.AnalysisConfig) []entity.Miner {
	miners := make([]entity.Miner, 0)
	for _, name := range config.Miners {
		factory, err := config.MinerAlgorithmFactory.Get(name)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to get mining factory for %s", name))
			continue
		}

		miner, err := factory.Make()
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to make mining algorithm for %s", name))
			continue
		}

		miners = append(miners, miner)
	}

	return miners
}

func buildSplittersFromInputAndMiningResults(config *entity.AnalysisConfig, miningResults map[string]entity.Miner) []entity.Splitter {
	splitters := make([]entity.Splitter, 0)
	for _, name := range config.Splitters {
		factory, err := config.SplittingAlgorithmFactory.Get(name)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to get splitting factory for %s", name))
			continue
		}

		splitter, err := factory.Make(miningResults)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to make splitting algorithm for %s", name))
			continue
		}

		splitters = append(splitters, splitter)
	}

	return splitters
}

func buildExpandersFromInputAndMiningResults(config *entity.AnalysisConfig, miningResults map[string]entity.Miner) []entity.Expander {
	expanders := make([]entity.Expander, 0)
	for _, name := range config.Expanders {
		factory, err := config.ExpansionAlgorithmFactory.Get(name)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to get expasions factory for %s", name))
			continue
		}

		expander, err := factory.Make(miningResults)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to make expansion algorithm for %s", name))
			continue
		}

		expanders = append(expanders, expander)
	}

	return expanders
}
