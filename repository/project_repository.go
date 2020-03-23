package repository

import (
	"context"
	"errors"
)

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

type ProjectRepository interface {
	Add(ctx context.Context, project Project) error
	GetByURL(ctx context.Context, url string) (Project, error)
}
