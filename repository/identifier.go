package repository

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
)

var (
	// ErrIdentifierNoResults indicates that no identifiers were found matching the given criteria.
	ErrIdentifierNoResults = errors.New("no identifiers found for the given criteria")
	// ErrIdentifierUnexpected indicates that an error occurred while trying to perform an operation on IdentifierRepository.
	ErrIdentifierUnexpected = errors.New("unexpected error performing the current operation on IdentifierRepository")
)

// IdentifierRepository represents a repository able to store and retrieve identifiers.
type IdentifierRepository interface {
	// Add associates an identifier with a given Analysis.
	Add(ctx context.Context, analysis entity.AnalysisResults, ident entity.Identifier) error
	// FindAllByProject retrives a list of identifiers associated to the given project reference.
	FindAllByProject(ctx context.Context, projectRef string) ([]entity.Identifier, error)
	// FindAllByProjectAndFile retrieve a list of identifiers that match the given criteria.
	FindAllByProjectAndFile(ctx context.Context, projectRef string, filename string) ([]entity.Identifier, error)
}
