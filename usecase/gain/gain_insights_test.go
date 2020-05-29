package gain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase/gain"
	"github.com/stretchr/testify/assert"
)

func TestNewGainInsightsUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := gain.NewGainInsightsUsecase(nil, nil)

	assert.NotNil(t, uc)
}

func TestProcess_OnGainInsightsUsecase_WhenNoIdentifiers_ShouldReturnError(t *testing.T) {
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{},
		err:    repository.ErrIdentifierNoResults,
	}

	uc := gain.NewGainInsightsUsecase(identifierRepositoryMock, nil)

	insights, err := uc.Process(context.TODO(), "eroatta/test")

	assert.EqualError(t, err, gain.ErrIdentifiersNotFound.Error())
	assert.Empty(t, insights)
}

func TestProcess_OnGainInsightsUsecase_WhenErrorReadingIdentifiers_ShouldReturnError(t *testing.T) {
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{},
		err:    repository.ErrIdentifierUnexpected,
	}

	uc := gain.NewGainInsightsUsecase(identifierRepositoryMock, nil)

	insights, err := uc.Process(context.TODO(), "eroatta/test")

	assert.EqualError(t, err, gain.ErrUnableToReadIdentifiers.Error())
	assert.Empty(t, insights)
}

func TestProcess_OnGainInsightsUsecase_WhenFailingToSaveInsights_ShouldReturnError(t *testing.T) {
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{
			{Package: "main", File: "main.go"},
		},
		err: nil,
	}

	insightsRepositoryMock := insightsRepositoryMock{
		err: repository.ErrInsightUnexpected,
	}

	uc := gain.NewGainInsightsUsecase(identifierRepositoryMock, insightsRepositoryMock)

	insights, err := uc.Process(context.TODO(), "eroatta/test")

	assert.EqualError(t, err, gain.ErrUnableToGainInsights.Error())
	assert.Empty(t, insights)
}

func TestProcess_OnGainInsightsUsecase_ShouldReturnInsightsByPackage(t *testing.T) {
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{
			{Package: "main", File: "main.go"},
			{Package: "main", File: "main.go"},
			{Package: "main", File: "helper.go"},
			{Package: "main_test", File: "main_test.go"},
		},
		err: nil,
	}

	insightsRepositoryMock := insightsRepositoryMock{
		err: nil,
	}

	uc := gain.NewGainInsightsUsecase(identifierRepositoryMock, insightsRepositoryMock)

	insights, err := uc.Process(context.TODO(), "eroatta/test")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(insights))
	assert.ElementsMatch(t, []entity.Insight{
		{
			ProjectRef: "eroatta/test",
			Package:    "main",
		},
		{
			ProjectRef: "eroatta/test",
			Package:    "main_test",
		},
	}, insights)
}

type identifierRepositoryMock struct {
	idents []entity.Identifier
	err    error
}

func (i identifierRepositoryMock) Add(ctx context.Context, project entity.Project, ident entity.Identifier) error {
	return errors.New("shouldn't be called")
}

func (i identifierRepositoryMock) FindAllByProject(ctx context.Context, projectRef string) ([]entity.Identifier, error) {
	return i.idents, i.err
}

type insightsRepositoryMock struct {
	err error
}

func (i insightsRepositoryMock) AddAll(ctx context.Context, insights []entity.Insight) error {
	return i.err
}
