package repository

import (
	"context"

	"github.com/eroatta/src-reader/entity"
)

// RemoteProjectRepository represents a repository capable of retrieving metadata from a remote project.
type RemoteProjectRepository interface {
	// RetrieveMetada extracts context information from a remote repository.
	RetrieveMetadata(ctx context.Context, remoteRepository string) (entity.Metadata, error)
}
