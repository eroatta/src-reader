package usecase

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/repository"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// ErrAnalysisNotFound indicates there is no analysis to be processed on the usecase.
var ErrAnalysisNotFound = errors.New("analysis not found")

// DeleteAnalysisUsecase defines the contract for the use case related to delete an existing analysis.
type DeleteAnalysisUsecase interface {
	// Process handles the process to delete an existing Analysis, given its ID.
	Process(ctx context.Context, analysisID uuid.UUID) error
}

// NewDeleteAnalysisUsecase create initializes a DeleteAnalysisUsecase instance.
func NewDeleteAnalysisUsecase(duc DeleteInsightsUsecase, ir repository.IdentifierRepository, ar repository.AnalysisRepository) DeleteAnalysisUsecase {
	return deleteAnalysisUsecase{
		deleteInsightsUsecase: duc,
		ir:                    ir,
		ar:                    ar,
	}
}

type deleteAnalysisUsecase struct {
	deleteInsightsUsecase DeleteInsightsUsecase
	ir                    repository.IdentifierRepository
	ar                    repository.AnalysisRepository
}

func (uc deleteAnalysisUsecase) Process(ctx context.Context, analysisID uuid.UUID) error {
	err := uc.deleteInsightsUsecase.Process(ctx, analysisID)
	if err == ErrUnexpected {
		log.Errorf("unable to execute delete insights usecase for analysis ID: %v", analysisID)
		return ErrUnexpected
	}

	err = uc.ir.DeleteAllByAnalysisID(ctx, analysisID)
	if err == repository.ErrIdentifierUnexpected {
		log.Errorf("unable to delete identifiers for analysis ID: %v", analysisID)
		return ErrUnexpected
	}

	err = uc.ar.Delete(ctx, analysisID)
	switch err {
	case nil:
		// do nothing
	case repository.ErrAnalysisNoResults:
		return ErrAnalysisNotFound
	default:
		log.Errorf("unable to delete analysis with ID: %v", analysisID)
		return ErrUnexpected
	}

	return nil
}
