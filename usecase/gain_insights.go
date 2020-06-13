package usecase

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrIdentifiersNotFound indicates that no identifiers were found to process.
	ErrIdentifiersNotFound = errors.New("no identifiers found for the project")
	// ErrUnableToReadIdentifiers indicates that the identifiers couldn't be retrieved.
	ErrUnableToReadIdentifiers = errors.New("unable to retrieve and read identifiers")
	// ErrUnableToSaveInsights indicates that an error occurred while trying to store the insights.
	ErrUnableToGainInsights = errors.New("unable to gain insights from identifiers")
)

type GainInsightsUsecase interface {
	Process(ctx context.Context, projectRef string) ([]entity.Insight, error)
}

func NewGainInsightsUsecase(identr repository.IdentifierRepository, insr repository.InsightRepository) GainInsightsUsecase {
	return gainInsightsUsecase{
		identr: identr,
		insr:   insr,
	}
}

type gainInsightsUsecase struct {
	identr repository.IdentifierRepository
	insr   repository.InsightRepository
}

func (uc gainInsightsUsecase) Process(ctx context.Context, projectRef string) ([]entity.Insight, error) {
	// grab each identifier
	identifiers, err := uc.identr.FindAllByProject(ctx, projectRef)
	switch err {
	case nil:
		// do nothing
	case repository.ErrIdentifierNoResults:
		return []entity.Insight{}, ErrIdentifiersNotFound
	default:
		log.WithError(err).Errorf("unable to retrieve identifiers for %s", projectRef)
		return []entity.Insight{}, ErrUnableToReadIdentifiers
	}

	byPackages := make(map[string]entity.Insight)
	for _, ident := range identifiers {
		metrics, ok := byPackages[ident.FullPackageName()]
		if !ok {
			metrics = entity.Insight{
				ProjectRef:      projectRef,
				Package:         ident.Package,
				TotalSplits:     make(map[string]int),
				TotalExpansions: make(map[string]int),
				Files:           make(map[string]struct{}),
			}
		}

		metrics.Include(ident)
		byPackages[ident.FullPackageName()] = metrics
	}

	insights := asArray(byPackages)
	err = uc.insr.AddAll(ctx, insights)
	if err != nil {
		log.WithError(err).Errorf("unable to store insights for %s", projectRef)
		return []entity.Insight{}, ErrUnableToGainInsights
	}

	return insights, nil
}

func asArray(input map[string]entity.Insight) []entity.Insight {
	insights := make([]entity.Insight, 0)
	for _, v := range input {
		insights = append(insights, v)
	}

	return insights
}
