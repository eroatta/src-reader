package repository

import (
	"context"
	"errors"
)

// TODO: each repository should handle their own errors
var (
	// ErrNoResults indicates that no results were found.
	ErrNoResults = errors.New("No results found")
	// ErrUnexpected indicates that the current action couldn't be completed because of an internal issue.
	ErrUnexpected = errors.New("Unexpected error performing the current action")
)

type Project struct {
	Status   string
	Metadata Metadata
}

// ProjectRepository represents a repository capable of operating with new and existing projects.
type ProjectRepository interface {
	// Add adds a new Project to the current repository.
	Add(ctx context.Context, project Project) error
	// GetByURL retrieves a Project using the remote URL.
	GetByURL(ctx context.Context, url string) (Project, error)
}
