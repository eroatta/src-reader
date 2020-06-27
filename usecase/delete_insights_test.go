package usecase_test

import (
	"context"
	"testing"

	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDeleteInsightsUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewDeleteInsightsUsecase(nil)

	assert.Empty(t, uc)
}

func TestProcess_OnDeleteInsightsUsecase_WhenNoExistingInsights_ShouldReturnInsightsNotFound(t *testing.T) {
	insightsRepositoryMock := insightsRepositoryMock{
		delErr: repository.ErrInsightNoResults,
	}
	uc := usecase.NewDeleteInsightsUsecase(insightsRepositoryMock)

	ID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), ID)

	assert.EqualError(t, err, usecase.ErrInsightsNotFound.Error())
}

func TestProcess_OnDeleteInsightsUsecase_WhenErrorDeletingInsights_ShouldReturnError(t *testing.T) {
	insightsRepositoryMock := insightsRepositoryMock{
		delErr: repository.ErrInsightUnexpected,
	}
	uc := usecase.NewDeleteInsightsUsecase(insightsRepositoryMock)

	ID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), ID)

	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnDeleteInsightsUsecase_WhenInsightsDeleted_ShouldReturnNoError(t *testing.T) {
	insightsRepositoryMock := insightsRepositoryMock{
		delErr: nil,
	}
	uc := usecase.NewDeleteInsightsUsecase(insightsRepositoryMock)

	ID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), ID)

	assert.NoError(t, err)
}
