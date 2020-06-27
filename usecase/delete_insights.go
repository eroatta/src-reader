package usecase

import (
	"context"

	"github.com/eroatta/src-reader/repository"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type DeleteInsightsUsecase interface {
	Process(ctx context.Context, analsysisID uuid.UUID) error
}

func NewDeleteInsightsUsecase(ir repository.InsightRepository) DeleteInsightsUsecase {
	return deleteInsightsUsecase{
		insightsRepository: ir,
	}
}

type deleteInsightsUsecase struct {
	insightsRepository repository.InsightRepository
}

func (uc deleteInsightsUsecase) Process(ctx context.Context, analysisID uuid.UUID) error {
	err := uc.insightsRepository.DeleteAllByAnalysisID(ctx, analysisID)
	switch err {
	case nil:
		// do nothing
	case repository.ErrInsightNoResults:
		return ErrInsightsNotFound
	default:
		log.WithError(err).Errorf("unable to delete insights for analysis ID: %v", analysisID)
		return ErrUnexpected
	}

	return nil
}
