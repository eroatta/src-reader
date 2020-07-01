package repository

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
	"github.com/google/uuid"
)

var (
	// ErrProjectNoResults indicates that no projects were found matching the given criteria.
	ErrProjectNoResults = errors.New("no projects found for the given criteria")
	// ErrProjectUnexpected indicates that the current action couldn't be completed because of an internal issue.
	ErrProjectUnexpected = errors.New("unexpected error performing the current action")
)

// ProjectRepository represents a repository capable of operating with new and existing projects.
type ProjectRepository interface {
	// Add adds a new Project to the current repository.
	Add(ctx context.Context, project entity.Project) error
	// Get retrieves a Project by ID.
	Get(ctx context.Context, ID uuid.UUID) (entity.Project, error)
	// GetByReference retrieves a Project using its reference name.
	GetByReference(ctx context.Context, projectRef string) (entity.Project, error)
	// Delete removes an existing Project from the current repository.
	Delete(ctx context.Context, ID uuid.UUID) error
}
