package repository

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
)

var (
	ErrSourceCodeUnableCreateDestination = errors.New("unable to create destination for source code")
	ErrSourceCodeCloneRemoteRepository   = errors.New("unable to clone remote repository")
	ErrSourceCodeUnableAccessMetadata    = errors.New("unable to access source code metadata information")
)

// SourceCodeRepository represents a repository capable of handle source code.
type SourceCodeRepository interface {
	// Clone clones the source code, under a given name, using the provided clone URL.
	Clone(ctx context.Context, fullname string, cloneURL string) (entity.SourceCode, error)
}
