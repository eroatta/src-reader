package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
)

var (
	ErrUnableToReadProject      = errors.New("")
	ErrUnableToRetrieveMetadata = errors.New("")
	ErrUnableToCloneSourceCode  = errors.New("")
	ErrUnableToSaveProject      = errors.New("")
)

// ImportProjectUsecase defines the contract for the use case related to the import process
// for a project.
type ImportProjectUsecase interface {
	// Import executes the pipeline to import a project from GitHub and
	// extract the identifiers. It returns the project informartion.
	Import(ctx context.Context, url string) (repository.Project, error)
}

// NewImportProjectUsecase initializes a new ImportProjectUsecase handler.
func NewImportProjectUsecase(pr repository.ProjectRepository, rpr repository.RemoteProjectRepository) ImportProjectUsecase {
	return importProjectUsecase{
		projectRepository:       pr,
		remoteProjectRepository: rpr,
	}
}

type importProjectUsecase struct {
	projectRepository       repository.ProjectRepository
	remoteProjectRepository repository.RemoteProjectRepository
	sourceCodeRepository    repository.SourceCodeRepository
}

func (uc importProjectUsecase) Import(ctx context.Context, url string) (repository.Project, error) {
	// check if not previously imported
	project, err := uc.projectRepository.GetByURL(ctx, url)
	switch err {
	case nil:
		return project, nil
	case repository.ErrNoResults:
		// continue
	default:
		log.WithError(err).Error(fmt.Sprintf("unable to retrieve project for %s", url))
		return repository.Project{}, ErrUnableToReadProject
	}

	// retrieve metadata
	metadata, err := uc.remoteProjectRepository.RetrieveMetadata(ctx, url)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("unable to retrive metadata for %s", url))
		return repository.Project{}, ErrUnableToRetrieveMetadata
	}

	project = repository.Project{
		Status:   "retrieved",
		Metadata: metadata,
	}

	// clone the source code
	err = uc.sourceCodeRepository.Clone(ctx, url)
	if err != nil {
		// TODO: handle error
	}

	// store the results
	err = uc.projectRepository.Add(ctx, project)
	if err != nil {
		// TODO: handle error
	}

	return repository.Project{}, nil
}
