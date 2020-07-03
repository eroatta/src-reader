package usecase

import (
	"context"

	"github.com/eroatta/src-reader/repository"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// DeleteInsightsUsecase defines the contract for the use case related to delete insights extracted
// from an existing analysis.
type DeleteInsightsUsecase interface {
	// Process handles the process to delete the Insights related to a given Analysis ID.
	Process(ctx context.Context, analsysisID uuid.UUID) error
}

// NewDeleteInsightsUsecase create initializes a DeleteInsightsUsecase instance.
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
