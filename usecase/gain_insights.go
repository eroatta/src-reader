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
	// ErrPreviousInsightsFound indicates there are existing insights.
	ErrPreviousInsightsFound = errors.New("existing previous insights for the analysis")
	// ErrIdentifiersNotFound indicates that no identifiers were found to process.
	ErrIdentifiersNotFound = errors.New("no identifiers found for the project")
	// ErrUnableToReadIdentifiers indicates that the identifiers couldn't be retrieved.
	ErrUnableToReadIdentifiers = errors.New("unable to retrieve and read identifiers")
	// ErrUnableToGainInsights indicates that an error occurred while trying to store the insights.
	ErrUnableToGainInsights = errors.New("unable to gain insights from identifiers")
)

// GainInsightsUsecase defines the contract for the usecase to analyze results and gain insights
// from a previous analysis.
type GainInsightsUsecase interface {
	Process(ctx context.Context, analysisID uuid.UUID) ([]entity.Insight, error)
}

// NewGainInsightsUsecase initializes a new GainInsightsUsecase instance.
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

func (uc gainInsightsUsecase) Process(ctx context.Context, analysisID uuid.UUID) ([]entity.Insight, error) {
	insights, err := uc.insr.GetByAnalysisID(ctx, analysisID)
	switch err {
	case repository.ErrInsightNoResults:
		// do nothing
	case nil:
		return insights, ErrPreviousInsightsFound
	default:
		log.WithError(err).Errorf("unable to check for previous insights on analysis ID: %v", analysisID)
		return []entity.Insight{}, ErrUnableToGainInsights
	}

	// grab each identifier
	identifiers, err := uc.identr.FindAllByAnalysisID(ctx, analysisID)
	switch err {
	case nil:
		// do nothing
	case repository.ErrIdentifierNoResults:
		return []entity.Insight{}, ErrIdentifiersNotFound
	default:
		log.WithError(err).Errorf("unable to retrieve identifiers for analysis ID: %v", analysisID)
		return []entity.Insight{}, ErrUnableToReadIdentifiers
	}

	byPackages := make(map[string]entity.Insight)
	for _, ident := range identifiers {
		metrics, ok := byPackages[ident.FullPackageName()]
		if !ok {
			metrics = entity.Insight{
				ProjectRef:      ident.ProjectRef,
				AnalysisID:      analysisID,
				Package:         ident.Package,
				TotalSplits:     make(map[string]int),
				TotalExpansions: make(map[string]int),
				Files:           make(map[string]struct{}),
			}
		}

		metrics.Include(ident)
		byPackages[ident.FullPackageName()] = metrics
	}

	insights = asArray(byPackages)
	err = uc.insr.AddAll(ctx, insights)
	if err != nil {
		log.WithError(err).Errorf("unable to store insights for analysis ID: %v", analysisID)
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
