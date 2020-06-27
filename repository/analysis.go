package repository

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
	"github.com/google/uuid"
)

var (
	// ErrAnalysisNoResults indicates that no analysis were found matching the given criteria.
	ErrAnalysisNoResults = errors.New("no analysis found for the given criteria")
	// ErrAnalysisUnexpected indicates that an error occurred while trying to perform an operation on AnalysisRepository.
	ErrAnalysisUnexpected = errors.New("unexpected error performing the current operation on AnalysisRepository")
)

// AnalysisRepository represents a repository able to store and retrieve analysis results.
type AnalysisRepository interface {
	// Add adds a new Analysis Results to the current repository.
	Add(ctx context.Context, analysis entity.AnalysisResults) error
	// GetByProjectID retrieves an existing analysis for the given Project.
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (entity.AnalysisResults, error)
	// Delete removes an Analysis from the current repository, using its ID.
	Delete(ctx context.Context, id uuid.UUID) error
}
