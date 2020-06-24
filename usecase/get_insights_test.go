package usecase_test

import (
	"context"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewGetInsightsUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewGetInsightsUsecase(nil)

	assert.NotNil(t, uc)
}

func TestProcess_OnGetInsightsUsecase_WhenNoExistingInsights_ShouldReturnError(t *testing.T) {
	insightsRepositoryMock := insightsRepositoryMock{
		getErr: repository.ErrInsightNoResults,
	}

	uc := usecase.NewGetInsightsUsecase(insightsRepositoryMock)

	id, _ := uuid.NewUUID()
	insights, err := uc.Process(context.TODO(), id)

	assert.Empty(t, insights)
	assert.EqualError(t, err, usecase.ErrInsightsNotFound.Error())
}

func TestProcess_OnGetInsightsUsecase_WhenErrorRetrievingInsights_ShouldReturnError(t *testing.T) {
	insightsRepositoryMock := insightsRepositoryMock{
		getErr: repository.ErrInsightUnexpected,
	}

	uc := usecase.NewGetInsightsUsecase(insightsRepositoryMock)

	id, _ := uuid.NewUUID()
	insights, err := uc.Process(context.TODO(), id)

	assert.Empty(t, insights)
	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnGetInsightsUsecase_WhenExistingInsights_ShouldReturnInsights(t *testing.T) {
	analysisID, _ := uuid.NewUUID()
	insightsRepositoryMock := insightsRepositoryMock{
		insights: []entity.Insight{
			{
				ID:               "ed2cd46a-4afd-4d49-a6ea-1c8d12d40134",
				ProjectRef:       "test/mytest",
				AnalysisID:       analysisID,
				Package:          "main",
				TotalIdentifiers: 3,
				TotalExported:    1,
				TotalSplits: map[string]int{
					"conserv": 5,
				},
				TotalExpansions: map[string]int{
					"no_exp": 5,
				},
				TotalWeight: 2.267,
				Files: map[string]struct{}{
					"main.go":   {},
					"helper.go": {},
				},
			},
			{
				ID:               "ed2cd46a-4afd-4d49-a6ea-1c8d12d40135",
				ProjectRef:       "test/mytest",
				AnalysisID:       analysisID,
				Package:          "main_test",
				TotalIdentifiers: 1,
				TotalExported:    1,
				TotalSplits: map[string]int{
					"conserv": 1,
				},
				TotalExpansions: map[string]int{
					"no_exp": 1,
				},
				TotalWeight: 1.0,
				Files: map[string]struct{}{
					"main_test.go": {},
				},
			},
		},
	}

	uc := usecase.NewGetInsightsUsecase(insightsRepositoryMock)

	insights, err := uc.Process(context.TODO(), analysisID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(insights))
}
