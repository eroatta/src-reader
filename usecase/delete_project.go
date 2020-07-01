package usecase

import (
	"context"

	"github.com/eroatta/src-reader/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type DeleteProjectUsecase interface {
	Process(ctx context.Context, projectID uuid.UUID) error
}

func NewDeleteProjectUsecase(duc DeleteAnalysisUsecase, ar repository.AnalysisRepository,
	scr repository.SourceCodeRepository, pr repository.ProjectRepository) DeleteProjectUsecase {
	return deleteProjectUsecase{
		deleteAnalysisUsecase: duc,
		ar:                    ar,
		scr:                   scr,
		pr:                    pr,
	}
}

type deleteProjectUsecase struct {
	deleteAnalysisUsecase DeleteAnalysisUsecase
	ar                    repository.AnalysisRepository
	scr                   repository.SourceCodeRepository
	pr                    repository.ProjectRepository
}

func (uc deleteProjectUsecase) Process(ctx context.Context, projectID uuid.UUID) error {
	analysis, err := uc.ar.GetByProjectID(ctx, projectID)
	if err == repository.ErrAnalysisUnexpected {
		logrus.Errorf("unable to access analsys for project ID: %v", projectID)
		return ErrUnexpected
	}

	if err != nil && uc.deleteAnalysisUsecase.Process(ctx, analysis.ID) == ErrUnexpected {
		log.Errorf("unable to execute delete analysis usecase for project ID: %v", projectID)
		return ErrUnexpected
	}

	project, err := uc.pr.Get(ctx, projectID)
	if err != nil && err == repository.ErrProjectUnexpected {
		log.Errorf("unable to retrieve project ID: %v", projectID)
		return ErrUnexpected
	}

	if err != ErrProjectNotFound && uc.scr.Remove(ctx, project.SourceCode.Location) == repository.ErrSourceCodeUnableToRemove {
		log.Errorf("unable to remove source code for project ID: %v", projectID)
		return ErrUnexpected
	}

	switch uc.pr.Delete(ctx, projectID) {
	case nil:
		// do nothing
	case repository.ErrProjectNoResults:
		return ErrProjectNotFound
	default:
		log.Errorf("unable to delete project with ID: %v", projectID)
		return ErrUnexpected
	}

	return nil
}
