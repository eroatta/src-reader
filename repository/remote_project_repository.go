package repository

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
)

var (
	// ErrMetadataUnexpected indicates that the current action couldn't be completed because of an internal issue.
	ErrMetadataUnexpected = errors.New("unexpected error performing the current action")
)

// RemoteProjectRepository represents a repository capable of retrieving metadata from a remote project.
type RemoteProjectRepository interface {
	// RetrieveMetada extracts context information from a remote repository.
	RetrieveMetadata(ctx context.Context, remoteRepository string) (entity.Metadata, error)
}
