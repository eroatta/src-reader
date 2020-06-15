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

func TestNewGainInsightsUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewGainInsightsUsecase(nil, nil)

	assert.NotNil(t, uc)
}

func TestProcess_OnGainInsightsUsecase_WhenNoIdentifiers_ShouldReturnError(t *testing.T) {
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{},
		err:    repository.ErrIdentifierNoResults,
	}

	uc := usecase.NewGainInsightsUsecase(identifierRepositoryMock, nil)

	analysisID, _ := uuid.NewUUID()
	insights, err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrIdentifiersNotFound.Error())
	assert.Empty(t, insights)
}

func TestProcess_OnGainInsightsUsecase_WhenErrorReadingIdentifiers_ShouldReturnError(t *testing.T) {
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{},
		err:    repository.ErrIdentifierUnexpected,
	}

	uc := usecase.NewGainInsightsUsecase(identifierRepositoryMock, nil)

	analysisID, _ := uuid.NewUUID()
	insights, err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrUnableToReadIdentifiers.Error())
	assert.Empty(t, insights)
}

func TestProcess_OnGainInsightsUsecase_WhenFailingToSaveInsights_ShouldReturnError(t *testing.T) {
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{
			{Package: "main", File: "main.go", Name: "main"},
		},
		err: nil,
	}

	insightsRepositoryMock := insightsRepositoryMock{
		err: repository.ErrInsightUnexpected,
	}

	uc := usecase.NewGainInsightsUsecase(identifierRepositoryMock, insightsRepositoryMock)

	analysisID, _ := uuid.NewUUID()
	insights, err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrUnableToGainInsights.Error())
	assert.Empty(t, insights)
}

func TestProcess_OnGainInsightsUsecase_ShouldReturnInsightsByPackage(t *testing.T) {
	analysisID, _ := uuid.NewUUID()
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{
			{
				Package: "main",
				File:    "main.go",
				Name:    "main",
				Splits: map[string][]entity.Split{
					"conserv": {{Order: 1, Value: "main"}},
				},
				Expansions: map[string][]entity.Expansion{
					"no_exp": {{From: "main", Values: []string{"main"}}},
				},
				Normalization: entity.Normalization{
					Word:      "main",
					Algorithm: "conserv+no_exp",
					Score:     0.98,
				},
				AnalysisID: analysisID,
				ProjectRef: "test/mytest",
			},
			{
				Package: "main",
				File:    "main.go",
				Name:    "helperFunc",
				Splits: map[string][]entity.Split{
					"conserv": {
						{Order: 1, Value: "helper"},
						{Order: 2, Value: "func"},
					},
				},
				Expansions: map[string][]entity.Expansion{
					"no_exp": {
						{From: "main", Values: []string{"main"}},
						{From: "func", Values: []string{"function"}},
					},
				},
				Normalization: entity.Normalization{
					Word:      "helperFunction",
					Algorithm: "conserv+no_exp",
					Score:     0.93,
				},
				AnalysisID: analysisID,
				ProjectRef: "test/mytest",
			},
			{
				Package: "main",
				File:    "helper.go",
				Name:    "HelperFunc",
				Splits: map[string][]entity.Split{
					"conserv": {
						{Order: 1, Value: "helper"},
						{Order: 2, Value: "func"},
					},
				},
				Expansions: map[string][]entity.Expansion{
					"no_exp": {
						{From: "main", Values: []string{"main"}},
						{From: "func", Values: []string{"function"}},
					},
				},
				Normalization: entity.Normalization{
					Word:      "HelperFunction",
					Algorithm: "conserv+no_exp",
					Score:     0.93,
				},
				AnalysisID: analysisID,
				ProjectRef: "test/mytest",
			},
			{
				Package: "main_test",
				File:    "main_test.go",
				Name:    "Mock",
				Splits: map[string][]entity.Split{
					"conserv": {
						{Order: 1, Value: "mock"},
					},
				},
				Expansions: map[string][]entity.Expansion{
					"no_exp": {
						{From: "mock", Values: []string{"mock"}},
					},
				},
				Normalization: entity.Normalization{
					Word:      "Mock",
					Algorithm: "conserv+no_exp",
					Score:     1.0,
				},
				AnalysisID: analysisID,
				ProjectRef: "test/mytest",
			},
		},
		err: nil,
	}

	insightsRepositoryMock := insightsRepositoryMock{
		err: nil,
	}

	uc := usecase.NewGainInsightsUsecase(identifierRepositoryMock, insightsRepositoryMock)

	insights, err := uc.Process(context.TODO(), analysisID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(insights))
	assert.ElementsMatch(t, []entity.Insight{
		{
			ProjectRef:       "test/mytest",
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
			ProjectRef:       "test/mytest",
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
	}, insights)
}
