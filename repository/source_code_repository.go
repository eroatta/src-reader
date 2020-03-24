package repository

import (
	"context"
)

// SourceCodeRepository represents a repository capable of handle source code.
type SourceCodeRepository interface {
	// Clone clones the source code for a given project URL.
	Clone(ctx context.Context, url string) error
}
