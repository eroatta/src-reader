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
	// ErrInErrInsightUnexpected indicates that an error occurred while trying to perform an operation on InsightRespository.
	ErrInsightUnexpected = errors.New("unexpected error performing the current operation on Insightrepository")
)

type InsightRepository interface {
	AddAll(ctx context.Context, insights []entity.Insight) error
	GetByAnalysisID(ctx context.Context, analysisID uuid.UUID) ([]entity.Insight, error)
}
