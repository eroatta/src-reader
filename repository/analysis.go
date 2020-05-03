package repository

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
)

var (
	// ErrAnalysisUnexpected indicates that an error occurred while trying to perform an operation on AnalysisRepository.
	ErrAnalysisUnexpected = errors.New("unexpected error performing the current operation on AnalysisRepository")
)

// AnalysisRepository represents a repository able to store and retrieve analysis results.
type AnalysisRepository interface {
	// Add adds a new Analysis Results to the current repository.
	Add(ctx context.Context, analysis entity.AnalysisResults) error
}
