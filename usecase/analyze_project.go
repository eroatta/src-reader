package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
)

var (
	ErrUnableToRetrieveFiles = errors.New("unable to retrieve files for processing")
	ErrUnableToBuildASTs     = errors.New("unable to create ASTs from input")
	ErrUnableToMineASTs      = errors.New("unable to apply one or more miners to the ASTs")
)

type AnalyzeProjectUsecase interface {
	Analyze(ctx context.Context, project entity.Project) (Results, error)
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

func (uc analyzeProjectUsecase) Analyze(ctx context.Context, project entity.Project) (Results, error) {
	// read filenames from project and get the files
	loc := project.SourceCode.Location
	for _, file := range project.SourceCode.Files {
		_, err := uc.sourceCodeRepository.Read(ctx, loc, file)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("unable to read file %s/%s from source code repository", loc, file))
			return Results{}, ErrUnableToRetrieveFiles
		}
	}

	// parse them and retrieve the AST
	// apply the pre-process step (mine them)
	// analyze each identifier
	// return the splits and expansions found
	// store them
	return Results{}, nil
}
