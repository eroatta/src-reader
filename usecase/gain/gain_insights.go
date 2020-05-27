package gain

import (
	"context"

	"github.com/eroatta/src-reader/repository"
)

type GainInsightsUsecase interface {
	Process(ctx context.Context, projectRef string) error
}

type gainInsightsUsecase struct {
	ir repository.IdentifierRepository
}

func (guc gainInsightsUsecase) Process(ctx context.Context, projectRef string) error {
	// grab each identifier
	// calculate metric (dev standard) for each algorithm
	// pick the best one
	// update each identifier o store each one as a new revision
	// group the results by package
	// store the results
	return nil
}
