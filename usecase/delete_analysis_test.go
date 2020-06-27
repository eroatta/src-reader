package usecase_test

import (
	"context"
	"testing"

	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDeleteAnalysisUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewDeleteAnalysisUsecase(nil, nil, nil)

	assert.Empty(t, uc)
}

func TestProcess_OnDeleteAnalysisUsecase_WhenErrorDeletingInsights_ShouldReturnError(t *testing.T) {
	duc := deleteInsightsUsecaseMock{
		err: usecase.ErrUnexpected,
	}
	uc := usecase.NewDeleteAnalysisUsecase(duc, nil, nil)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnDeleteAnalysisUsecase_WhenErrorDeletingIdentifiers_ShouldReturnError(t *testing.T) {
	duc := deleteInsightsUsecaseMock{
		err: nil,
	}
	ir := identifierRepositoryMock{
		delErr: repository.ErrIdentifierUnexpected,
	}
	uc := usecase.NewDeleteAnalysisUsecase(duc, ir, nil)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnDeleteAnalysisUsecase_WhenErrorDeletingAnalysis_ShouldReturnError(t *testing.T) {
	duc := deleteInsightsUsecaseMock{
		err: nil,
	}
	ir := identifierRepositoryMock{
		delErr: nil,
	}
	ar := analysisRepositoryMock{
		delErr: repository.ErrAnalysisUnexpected,
	}
	uc := usecase.NewDeleteAnalysisUsecase(duc, ir, ar)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnDeleteAnalysisUsecase_WhenEverythingDeleted_ShouldReturnNoError(t *testing.T) {
	duc := deleteInsightsUsecaseMock{
		err: nil,
	}
	ir := identifierRepositoryMock{
		delErr: nil,
	}
	ar := analysisRepositoryMock{
		delErr: nil,
	}
	uc := usecase.NewDeleteAnalysisUsecase(duc, ir, ar)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.NoError(t, err)
}

func TestProcess_OnDeleteAnalysisUsecase_WhenNoInsights_ShouldReturnNoError(t *testing.T) {
	duc := deleteInsightsUsecaseMock{
		err: usecase.ErrInsightsNotFound,
	}
	ir := identifierRepositoryMock{
		delErr: nil,
	}
	ar := analysisRepositoryMock{
		delErr: nil,
	}
	uc := usecase.NewDeleteAnalysisUsecase(duc, ir, ar)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.NoError(t, err)
}

func TestProcess_OnDeleteAnalysisUsecase_WhenIdentifiersNotFound_ShouldReturnNoError(t *testing.T) {
	duc := deleteInsightsUsecaseMock{
		err: usecase.ErrInsightsNotFound,
	}
	ir := identifierRepositoryMock{
		delErr: repository.ErrIdentifierNoResults,
	}
	ar := analysisRepositoryMock{
		delErr: nil,
	}
	uc := usecase.NewDeleteAnalysisUsecase(duc, ir, ar)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.NoError(t, err)
}

func TestProcess_OnDeleteAnalysisUsecase_WhenAnalysisNotFound_ShouldReturnError(t *testing.T) {
	duc := deleteInsightsUsecaseMock{
		err: usecase.ErrInsightsNotFound,
	}
	ir := identifierRepositoryMock{
		delErr: repository.ErrIdentifierNoResults,
	}
	ar := analysisRepositoryMock{
		delErr: repository.ErrAnalysisNoResults,
	}
	uc := usecase.NewDeleteAnalysisUsecase(duc, ir, ar)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrAnalysisNotFound.Error())
}

type deleteInsightsUsecaseMock struct {
	err error
}

func (m deleteInsightsUsecaseMock) Process(ctx context.Context, analysisID uuid.UUID) error {
	return m.err
}
