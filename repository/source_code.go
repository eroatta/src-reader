package repository

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
)

var (
	// ErrSourceCodeUnableCreateDestination indicates that the underlying cloned source code destination couldn't be created.
	ErrSourceCodeUnableCreateDestination = errors.New("unable to create destination for source code")
	// ErrSourceCodeUnableCloneRemoteRepository indicates the remote repository couldn't be reached.
	ErrSourceCodeUnableCloneRemoteRepository = errors.New("unable to clone remote repository")
	// ErrSourceCodeUnableAccessMetadata indicates that the metadata for the current remote repository couldn't be reached.
	ErrSourceCodeUnableAccessMetadata = errors.New("unable to access source code metadata information")
	// ErrSourceCodeUnableToRemove indicates the source code couldn't be removed from the underlying storage.
	ErrSourceCodeUnableToRemove = errors.New("unable to remove source code")
	// ErrSourceCodeUnableReadFile indicates that the requested file couldn't be accessed or read from the underlying storage.
	ErrSourceCodeUnableReadFile = errors.New("unable to access or read file")
	// ErrSourceCodeNotFound indicates the source code is not present on the underlying storage.
	ErrSourceCodeNotFound = errors.New("unable to locate source code")
)

// SourceCodeRepository represents a repository capable of handle source code.
type SourceCodeRepository interface {
	// Clone clones the source code, under a given name, using the provided clone URL.
	Clone(ctx context.Context, fullname string, cloneURL string) (entity.SourceCode, error)
	// Remove removes the source code on the given location.
	Remove(ctx context.Context, location string) error
	// Read reads the content of the given file, relative to the provided location.
	Read(ctx context.Context, location string, filename string) ([]byte, error)
}
