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
	ErrUnableToBuildASTs = errors.New("unable to create ASTs from input")
	ErrUnableToMineASTs  = errors.New("unable to apply one or more miners to the ASTs")
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

type processConfig struct {
	reader string
	parser string
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
	// TODO: remove
	fmt.Println(miningResults)

	identc := step.Extract(valid, config.ExtractorFactory)
	splittedc := step.Split(identc, []entity.Splitter{}...)
	expandedc := step.Expand(splittedc, []entity.Expander{}...)
	for i := range expandedc {
		log.Info(i)
	}

	// analyze each identifier
	// return the splits and expansions found
	// store them
	return Results{}, nil
}
