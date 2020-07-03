package usecase

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrFileNotFound indicates the requested file was not found.
	ErrFileNotFound = errors.New("requested file not found")
)

// OriginalFileUsecase defines the contract to retrieve a file with its original content.
type OriginalFileUsecase interface {
	Process(ctx context.Context, projectRef string, filename string) ([]byte, error)
}

// NewOriginalFileUsecase initializes a new OriginalFileUsecase instance.
func NewOriginalFileUsecase(pr repository.ProjectRepository, scr repository.SourceCodeRepository) OriginalFileUsecase {
	return &originalFileUsecase{
		pr:  pr,
		scr: scr,
	}
}

type originalFileUsecase struct {
	pr  repository.ProjectRepository
	scr repository.SourceCodeRepository
}

func (uc *originalFileUsecase) Process(ctx context.Context, projectRef string, filename string) ([]byte, error) {
	project, err := uc.pr.GetByReference(ctx, projectRef)
	switch err {
	case nil:
		// do nothing
	case repository.ErrProjectNoResults:
		return nil, ErrProjectNotFound
	default:
		log.WithError(err).Errorf("unable to retrieve project %s", projectRef)
		return nil, ErrUnexpected
	}

	found := false
	for _, f := range project.SourceCode.Files {
		if f == filename {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrFileNotFound
	}

	raw, err := uc.scr.Read(ctx, project.SourceCode.Location, filename)
	if err != nil {
		log.WithError(err).Errorf("unable to read file %s on project %s", filename, projectRef)
		return nil, ErrUnexpected
	}

	return raw, nil
}
