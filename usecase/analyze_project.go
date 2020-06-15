package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase/step"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrProjectNotFound indicates that the requested project is not accessible.
	ErrProjectNotFound = errors.New("unable to retrieve requested Project")
	// ErrPreviousAnalysisFound indicates there is an existing previous analysis.
	ErrPreviousAnalysisFound = errors.New("existing previous analysis for project")
	// ErrUnableToBuildASTs indicates that an error occurred while trying to read or parse the source code files to
	// build the required Abstract Syntax Trees.
	ErrUnableToBuildASTs = errors.New("unable to create ASTs from input")
	// ErrUnableToMineASTs indicates that an error occurred while trying to create or apply the miners specified
	// during the import process.
	ErrUnableToMineASTs = errors.New("unable to apply one or more miners to the ASTs")
	// ErrUnableToCreateProcessors indicates that an error occurred while trying to create or apply the splitters or
	// expanders during the import process.
	ErrUnableToCreateProcessors = errors.New("unable to create splitting or expansion algorithms")
	// ErrUnableToSaveIdentifiers indicates that an error occurred while trying to save an already processed identifier.
	ErrUnableToSaveIdentifiers = errors.New("unable to save extracted and processed indentifiers")
	// ErrUnableToSaveAnalysis indicates that an error occurred while trying to store the results for an import process.
	ErrUnableToSaveAnalysis = errors.New("unable to save analysis results after completed processing")
	// ErrUnexpected indicates that an unexpected error ocurring while analyzing the project.
	ErrUnexpected = errors.New("unexpected error")
)

// AnalyzeProjectUsecase defines the contract for the use case related to the analysis process of a project.
type AnalyzeProjectUsecase interface {
	// Process performs the splitting and expansion process on the source code belonging to the Project.
	Process(ctx context.Context, projecID uuid.UUID) (entity.AnalysisResults, error)
}

// NewAnalyzeProjectUsecase initializes a new AnalyzeProjectUsecase handler.
func NewAnalyzeProjectUsecase(pr repository.ProjectRepository, scr repository.SourceCodeRepository,
	ir repository.IdentifierRepository, ar repository.AnalysisRepository, config *entity.AnalysisConfig) AnalyzeProjectUsecase {
	return &analyzeProjectUsecase{
		projectRepository:    pr,
		sourceCodeRepository: scr,
		identifierRepository: ir,
		analysisRepository:   ar,
		defaultConfig:        config,
	}
}

type analyzeProjectUsecase struct {
	projectRepository    repository.ProjectRepository
	sourceCodeRepository repository.SourceCodeRepository
	identifierRepository repository.IdentifierRepository
	analysisRepository   repository.AnalysisRepository
	defaultConfig        *entity.AnalysisConfig
}

// Process processes the given Project, based on the configuration provided by the AnalysisConfig.
// The process reads the source code, applies the given miners, splitters and expanders and then stores the results.
// It's the default implementation for the use case.
func (uc analyzeProjectUsecase) Process(ctx context.Context, projectID uuid.UUID) (entity.AnalysisResults, error) {
	project, err := uc.projectRepository.Get(ctx, projectID)
	switch err {
	case nil:
		// do nothing
	case repository.ErrProjectNoResults:
		return entity.AnalysisResults{}, ErrProjectNotFound
	default:
		log.WithError(err).Errorf("unable to retrieve project %s", projectID.String())
		return entity.AnalysisResults{}, ErrUnexpected
	}

	previousAnalysis, err := uc.analysisRepository.GetByProjectID(ctx, projectID)
	switch err {
	case repository.ErrAnalysisNoResults:
		// do nothing
	case nil:
		return previousAnalysis, ErrPreviousAnalysisFound
	default:
		log.WithError(err).Errorf("unable to check for previous analysis on project %s", projectID.String())
		return entity.AnalysisResults{}, ErrUnexpected
	}

	analysisID, _ := uuid.NewUUID()
	analysisResults := entity.AnalysisResults{
		ID:                analysisID,
		DateCreated:       time.Now(),
		ProjectID:         project.ID,
		ProjectName:       project.Metadata.Fullname,
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

			log.WithError(file.Error).Warnf("unable to read or parse file %s at %s", file.Name, project.SourceCode.Location)
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
		log.Errorf("unable to read or parse any file on %s for project %s",
			project.SourceCode.Location, project.Reference)
		return entity.AnalysisResults{}, ErrUnableToBuildASTs
	}

	// apply the pre-process step (mine them)
	miners := buildMiners(uc.defaultConfig)
	for _, miner := range miners {
		analysisResults.PipelineMiners = append(analysisResults.PipelineMiners, miner.Name())
	}

	miningResults := step.Mine(valid, miners...)

	// make the splitters from input and mining results
	splitters := buildSplittersFromMiningResults(uc.defaultConfig, miningResults)
	if len(splitters) == 0 {
		log.WithField("desired", uc.defaultConfig.Splitters).Error("unable to create any splitter")
		return entity.AnalysisResults{}, ErrUnableToCreateProcessors
	}
	for _, splitter := range splitters {
		analysisResults.PipelineSplitters = append(analysisResults.PipelineSplitters, splitter.Name())
	}

	// make the expanders from input and mining results
	expanders := buildExpandersFromMiningResults(uc.defaultConfig, miningResults)
	if len(expanders) == 0 {
		log.WithField("desired", uc.defaultConfig.Expanders).Error("unable to create any expander")
		return entity.AnalysisResults{}, ErrUnableToCreateProcessors
	}
	for _, expander := range expanders {
		analysisResults.PipelineExpanders = append(analysisResults.PipelineExpanders, expander.Name())
	}

	// analyze each identifier
	identc := step.Extract(valid, uc.defaultConfig.ExtractorFactory)
	splittedc := step.Split(identc, splitters...)
	expandedc := step.Expand(splittedc, expanders...)
	normalizedc := step.Normalize(expandedc)

	identErrorSamples := make([]string, 0)
	for ident := range normalizedc {
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
				ident.Name, ident.File, analysisResults.ProjectName))
			return entity.AnalysisResults{}, ErrUnableToSaveIdentifiers
		}
	}
	analysisResults.IdentifiersValid = analysisResults.IdentifiersTotal - analysisResults.IdentifiersError
	analysisResults.IdentifiersErrorSamples = identErrorSamples

	err = uc.analysisRepository.Add(ctx, analysisResults)
	if err != nil {
		log.WithError(err).Errorf("unable to save analysis results for project %s", project.Reference)
		return entity.AnalysisResults{}, ErrUnableToSaveAnalysis
	}

	return analysisResults, nil
}

// buildMiners initializes a set of miners, making them exclusives on current process.
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

// buildSplittersFromMiningResults initializes a set of splitters from the mining results, making them exclusives on current process.
func buildSplittersFromMiningResults(config *entity.AnalysisConfig, miningResults map[string]entity.Miner) []entity.Splitter {
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

// buildExpandersFromMiningResults initializes a set of expanders from the mining results, making them exclusives on current process.
func buildExpandersFromMiningResults(config *entity.AnalysisConfig, miningResults map[string]entity.Miner) []entity.Expander {
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
