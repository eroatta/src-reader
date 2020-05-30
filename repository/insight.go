package repository

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
)

// ErrInErrInsightUnexpected indicates that an error occurred while trying to perform an operation on InsightRespository.
var ErrInsightUnexpected = errors.New("unexpected error performing the current operation on Insightrepository")

type InsightRepository interface {
	AddAll(ctx context.Context, insights []entity.Insight) error
}
