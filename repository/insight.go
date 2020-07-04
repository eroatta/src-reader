package repository

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
	"github.com/google/uuid"
)

var (
	// ErrInsightNoResults indicates that no insights were found matching the given criteria.
	ErrInsightNoResults = errors.New("no insights found for the given criteria")
	// ErrInsightUnexpected indicates that an error occurred while trying to perform an operation on InsightRespository.
	ErrInsightUnexpected = errors.New("unexpected error performing the current operation on Insightrepository")
)

// InsightRepository represents a repository able to store, retrieve and delete insights.
type InsightRepository interface {
	// AddAll adds the provided insights to the repository.
	AddAll(ctx context.Context, insights []entity.Insight) error
	// GetByAnalysisID retrieves a list of Insight for a given analysis ID.
	GetByAnalysisID(ctx context.Context, analysisID uuid.UUID) ([]entity.Insight, error)
	// DeleteAllByAnalysisID removes all the Insights related to an analysis ID.
	DeleteAllByAnalysisID(ctx context.Context, analysisID uuid.UUID) error
}
