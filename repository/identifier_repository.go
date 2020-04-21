package repository

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
)

var (
	// ErrIdentifierUnexpected indicates that an error occurred while trying to perform an operation on IdentifierRepository.
	ErrIdentifierUnexpected = errors.New("unexpect error performing the current operation on IdentifierRepository")
)

type IdentifierRepository interface {
	Add(ctx context.Context, project entity.Project, ident entity.Identifier) error
}
