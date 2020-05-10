package create

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrUnableToReadProject indicates that an error occurred while trying to read the entity.Project
	// from any repository.
	ErrUnableToReadProject = errors.New("Unable to check against existing projects")
	// ErrUnableToRetrieveMetadata indicates that an error occurred while trying to retrieve Metadata information
	// related to the imported project.
	ErrUnableToRetrieveMetadata = errors.New("Unable to retrieve project metadata")
	// ErrUnableToCloneSourceCode indicates that an error occurred while trying to retrieve or store the source code
	// related to the imported project.
	ErrUnableToCloneSourceCode = errors.New("Unable to access or clone the source code")
	// ErrUnableToSaveProject indicates that an error occurred while trying to save the imported entity.Project.
	ErrUnableToSaveProject = errors.New("Unable to store project changes")
)

// ImportProjectUsecase defines the contract for the use case related to the import process
// for a project.
type ImportProjectUsecase interface {
	// Import executes the pipeline to import a project from GitHub and
	// extract the identifiers. It returns the project information.
	Import(ctx context.Context, url string) (entity.Project, error)
}

// NewImportProjectUsecase initializes a new ImportProjectUsecase handler.
func NewImportProjectUsecase(pr repository.ProjectRepository, mr repository.MetadataRepository,
	scr repository.SourceCodeRepository) ImportProjectUsecase {
	return importProjectUsecase{
		projectRepository:    pr,
		metadataRepository:   mr,
		sourceCodeRepository: scr,
	}
}

type importProjectUsecase struct {
	projectRepository    repository.ProjectRepository
	metadataRepository   repository.MetadataRepository
	sourceCodeRepository repository.SourceCodeRepository
}

// Import executes the pipeline to import a project from GitHub and extract the identifiers.
// It returns the project information. It's the current default implementation for the use case.
func (uc importProjectUsecase) Import(ctx context.Context, url string) (entity.Project, error) {
	// check if not previously imported
	project, err := uc.projectRepository.GetByURL(ctx, url)
	switch err {
	case nil:
		return project, nil
	case repository.ErrProjectNoResults:
		// continue
	default:
		log.WithError(err).Error(fmt.Sprintf("unable to retrieve project for %s", url))
		return entity.Project{}, ErrUnableToReadProject
	}

	// retrieve metadata
	metadata, err := uc.metadataRepository.RetrieveMetadata(ctx, url)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("unable to retrive metadata for %s", url))
		return entity.Project{}, ErrUnableToRetrieveMetadata
	}

	project = entity.Project{
		ID:       fmt.Sprintf("%x", md5.Sum([]byte(metadata.Fullname))),
		URL:      url,
		Metadata: metadata,
		Status:   "in_process",
	}

	// clone the source code
	sourceCode, err := uc.sourceCodeRepository.Clone(ctx, project.Metadata.Fullname, project.Metadata.CloneURL)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("unable to clone source code for %s", url))
		return entity.Project{}, ErrUnableToCloneSourceCode
	}
	project.SourceCode = sourceCode
	project.Status = "done"

	// store the results
	err = uc.projectRepository.Add(ctx, project)
	if err != nil {
		defer uc.sourceCodeRepository.Remove(ctx, sourceCode.Location)
		log.WithError(err).Error(fmt.Sprintf("unable to save project for %s", url))
		return entity.Project{}, ErrUnableToSaveProject
	}

	return project, nil
}
