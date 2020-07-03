package usecase

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrInsightsNotFound indicates that no Insights were found.
	ErrInsightsNotFound = errors.New("no insights found")
)

// GetInsightsUsecase handles the retrieval of Insights by its ID.
type GetInsightsUsecase interface {
	// Process retrieves insights and their related information.
	Process(ctx context.Context, ID uuid.UUID) ([]entity.Insight, error)
}

// NewGetInsightsUsecase initializes a new GetInsightsUsecase instance.
func NewGetInsightsUsecase(ir repository.InsightRepository) GetInsightsUsecase {
	return getInsightsUsecase{
		insightsRepository: ir,
	}
}

type getInsightsUsecase struct {
	insightsRepository repository.InsightRepository
}

func (uc getInsightsUsecase) Process(ctx context.Context, ID uuid.UUID) ([]entity.Insight, error) {
	insights, err := uc.insightsRepository.GetByAnalysisID(ctx, ID)
	switch err {
	case nil:
		// do nothing
	case repository.ErrInsightNoResults:
		return []entity.Insight{}, ErrInsightsNotFound
	default:
		log.WithError(err).Errorf("unable to retrieve insights with ID %v", ID)
		return []entity.Insight{}, ErrUnexpected
	}

	return insights, nil
}
