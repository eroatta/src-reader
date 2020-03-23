package repository

import (
	"context"
)

type SourceCodeRepository interface {
	Clone(ctx context.Context, url string) error
}
