package usecase

import (
	"context"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// GetProjectUsecase handles the retrieval of a Project by its ID.
type GetProjectUsecase interface {
	// Process retrieves a project and its related information.
	Process(ctx context.Context, ID uuid.UUID) (entity.Project, error)
}

// NewGetPNewGetProjectUsecase initializes a new GetProjectUsecase instance.
func NewGetProjectUsecase(pr repository.ProjectRepository) GetProjectUsecase {
	return getProjectUsecase{
		projectRepository: pr,
	}
}

type getProjectUsecase struct {
	projectRepository repository.ProjectRepository
}

func (uc getProjectUsecase) Process(ctx context.Context, ID uuid.UUID) (entity.Project, error) {
	project, err := uc.projectRepository.Get(ctx, ID)
	switch err {
	case nil:
		// do nothing
	case repository.ErrProjectNoResults:
		return entity.Project{}, ErrProjectNotFound
	default:
		log.WithError(err).Errorf("unable to retrieve project with ID %v", ID)
		return entity.Project{}, ErrUnexpected
	}

	return project, nil
}
