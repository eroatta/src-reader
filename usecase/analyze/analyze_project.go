package analyze

import (
	"context"
	"errors"
	"fmt"

	"github.com/eroatta/src-reader/code"
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
	Analyze(ctx context.Context, project entity.Project, config *entity.AnalysisConfig) (Results, error)
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

type Results struct {
}

func (uc analyzeProjectUsecase) Analyze(ctx context.Context, project entity.Project, config *entity.AnalysisConfig) (Results, error) {
	// read and parse files
	filesc := step.Read(ctx, uc.sourceCodeRepository, project.SourceCode.Location, project.SourceCode.Files)
	parsed := step.Parse(filesc)
	files := step.Merge(parsed)

	valid := make([]code.File, 0)
	for _, file := range files {
		if file.Error != nil {
			log.WithError(file.Error).Warn(fmt.Sprintf("unable to read or parse file %s at %s", file.Name, project.SourceCode.Location))
			continue
		}
		valid = append(valid, file)
	}

	// if every file can't be parsed, then fail
	if len(valid) == 0 {
		log.Error(fmt.Sprintf("unable to read or parse any file on %s for project %s",
			project.SourceCode.Location, project.URL))
		return Results{}, ErrUnableToBuildASTs
	}

	// apply the pre-process step (mine them)
	miningResults := step.Mine(valid, config.Miners...)

	// make the splitters from input and mining results
	splitters := make([]entity.Splitter, 0)
	for _, name := range config.Splitters {
		factory, err := config.SplittingAlgorithmFactory.Get(name)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to get splitting factory for %s", name))
			continue
		}

		splitter, err := factory.Make(config.StaticInputs, miningResults)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to make splitting algorithm for %s", name))
			continue
		}

		splitters = append(splitters, splitter)
	}
	if len(splitters) == 0 {
		log.WithField("desired", config.Splitters).Error("unable to create any splitter")
		return Results{}, ErrUnableToCreateProcessors
	}

	// make the expanders from input and mining results
	expanders := make([]entity.Expander, 0)
	for _, name := range config.Expanders {
		factory, err := config.ExpansionAlgorithmFactory.Get(name)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to get expasions factory for %s", name))
			continue
		}

		expander, err := factory.Make(config.StaticInputs, miningResults)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to make expansion algorithm for %s", name))
			continue
		}

		expanders = append(expanders, expander)
	}
	if len(expanders) == 0 {
		log.WithField("desired", config.Expanders).Error("unable to create any expander")
		return Results{}, ErrUnableToCreateProcessors
	}

	// analyze each identifier
	identc := step.Extract(valid, config.ExtractorFactory)
	splittedc := step.Split(identc, splitters...)
	expandedc := step.Expand(splittedc, expanders...)

	//errs := make([]error, 0)
	for ident := range expandedc {
		err := uc.identifierRepository.Add(ctx, project, ident)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to save identifier %s, on file %s for project %s",
				ident.Name, ident.File, project.URL))
			//errs = append(errs, err)
			return Results{}, ErrUnableToSaveIdentifiers
		}
	}

	return Results{}, nil
}
