package repository

import (
	"context"

	"github.com/eroatta/src-reader/entity"
)

type IdentifierRepository interface {
	Add(ctx context.Context, project entity.Project, ident entity.Identifier) error
}
