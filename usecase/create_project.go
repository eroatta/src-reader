package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/google/uuid"
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

// CreateProjectUsecase handles the creation and import or a Project.
type CreateProjectUsecase interface {
	// Process retrieves a project from GitHub and imports it.
	Process(ctx context.Context, projectRef string) (entity.Project, error)
}

// NewCreateProjectUsecase initializes a new CreateProjectUsecase instance.
func NewCreateProjectUsecase(pr repository.ProjectRepository, mr repository.MetadataRepository,
	scr repository.SourceCodeRepository) CreateProjectUsecase {
	return createProjectUsecase{
		projectRepository:    pr,
		metadataRepository:   mr,
		sourceCodeRepository: scr,
	}
}

type createProjectUsecase struct {
	projectRepository    repository.ProjectRepository
	metadataRepository   repository.MetadataRepository
	sourceCodeRepository repository.SourceCodeRepository
}

// Process executes the pipeline to import a project from GitHub. It returns the project information.
func (uc createProjectUsecase) Process(ctx context.Context, projectRef string) (entity.Project, error) {
	// check if not previously imported
	project, err := uc.projectRepository.GetByReference(ctx, projectRef)
	switch err {
	case nil:
		return project, nil
	case repository.ErrProjectNoResults:
		// continue
	default:
		log.WithError(err).Errorf("unable to retrieve project for %s", projectRef)
		return entity.Project{}, ErrUnableToReadProject
	}

	// retrieve metadata
	metadata, err := uc.metadataRepository.RetrieveMetadata(ctx, projectRef)
	if err != nil {
		log.WithError(err).Errorf("unable to retrive metadata for %s", projectRef)
		return entity.Project{}, ErrUnableToRetrieveMetadata
	}

	id, _ := uuid.NewUUID()
	project = entity.Project{
		ID:        id.String(),
		Reference: projectRef,
		CreatedAt: time.Now(),
		Metadata:  metadata,
		Status:    "in_process",
	}

	// clone the source code
	sourceCode, err := uc.sourceCodeRepository.Clone(ctx, projectRef, project.Metadata.CloneURL)
	if err != nil {
		log.WithError(err).Errorf("unable to clone source code for %s", projectRef)
		return entity.Project{}, ErrUnableToCloneSourceCode
	}
	project.SourceCode = sourceCode
	project.Status = "done"

	// store the results
	err = uc.projectRepository.Add(ctx, project)
	if err != nil {
		defer uc.sourceCodeRepository.Remove(ctx, sourceCode.Location)
		log.WithError(err).Errorf("unable to save project for %s", projectRef)
		return entity.Project{}, ErrUnableToSaveProject
	}

	return project, nil
}
