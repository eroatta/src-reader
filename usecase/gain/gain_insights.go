package gain

import (
	"context"
	"fmt"

	"github.com/eroatta/src-reader/repository"
)

type GainInsightsUsecase interface {
	Process(ctx context.Context, projectRef string) error
}

func NewGainInsightsUsecase(identr repository.IdentifierRepository, insr repository.InsightsRepository) GainInsightsUsecase {
	return gainInsightsUsecase{
		identr: identr,
		insr:   insr,
	}
}

type gainInsightsUsecase struct {
	identr repository.IdentifierRepository
	insr   repository.InsightsRepository
}

func (guc gainInsightsUsecase) Process(ctx context.Context, projectRef string) error {
	// grab each identifier
	identifiers, err := guc.identr.FindAllByProject(ctx, projectRef)
	if err != nil {
		// TODO
	}

	for _, ident := range identifiers {
		fmt.Println(ident)
	}
	// calculate metric (dev standard) for each algorithm
	// pick the best one
	// update each identifier o store each one as a new revision
	// group the results by package
	// store the results
	return nil
}
