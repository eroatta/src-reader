package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
)

var (
	ErrUnableToReadProject      = errors.New("Unable to check against existing projects")
	ErrUnableToRetrieveMetadata = errors.New("Unable to retrieve project metadata")
	ErrUnableToCloneSourceCode  = errors.New("Unable to access or clone the source code")
	ErrUnableToSaveProject      = errors.New("Unable to store project changes")
)

// ImportProjectUsecase defines the contract for the use case related to the import process
// for a project.
type ImportProjectUsecase interface {
	// Import executes the pipeline to import a project from GitHub and
	// extract the identifiers. It returns the project informartion.
	Import(ctx context.Context, url string) (repository.Project, error)
}

// NewImportProjectUsecase initializes a new ImportProjectUsecase handler.
func NewImportProjectUsecase(pr repository.ProjectRepository, rpr repository.RemoteProjectRepository,
	scr repository.SourceCodeRepository) ImportProjectUsecase {
	return importProjectUsecase{
		projectRepository:       pr,
		remoteProjectRepository: rpr,
		sourceCodeRepository:    scr,
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
		Metadata: metadata,
	}

	// clone the source code
	err = uc.sourceCodeRepository.Clone(ctx, url)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("unable to clone source code for %s", url))
		return repository.Project{}, ErrUnableToCloneSourceCode
	}

	// store the results
	err = uc.projectRepository.Add(ctx, project)
	if err != nil {
		// TODO: should we delete the cloned project?
		log.WithError(err).Error(fmt.Sprintf("unable to save project for %s", url))
		return repository.Project{}, ErrUnableToSaveProject
	}

	// after every step is completed, the import process is done
	project.Status = "done"

	return project, nil
}
