package usecase_test

import (
	"context"
	"testing"

	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewImportProjectUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewImportProjectUsecase(nil, nil)

	assert.NotNil(t, uc)
}

func TestImport_OnImportProjectUsecase_ShouldReturnImportResults(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestImport_OnImportProjectUsecase_WhenAlreadyImportedProject_ShouldImportResults(t *testing.T) {
	prMock := projectRepositoryMock{
		project: repository.Project{
			Status: "finished",
		},
		err: nil,
	}
	uc := usecase.NewImportProjectUsecase(prMock, nil)

	project, err := uc.Import(context.TODO(), "https://github.com/test/mytest")

	assert.NoError(t, err)
	assert.NotEmpty(t, project)
	assert.Equal(t, "finished", project.Status)
}

func TestImport_OnImportProjectUsecase_WhenUnableToCheckExistingProject_ShouldReturnError(t *testing.T) {
	prMock := projectRepositoryMock{
		project: repository.Project{},
		err:     repository.ErrUnexpected,
	}
	uc := usecase.NewImportProjectUsecase(prMock, nil)

	project, err := uc.Import(context.TODO(), "https://github.com/test/mytest")

	assert.EqualError(t, err, usecase.ErrUnableToReadProject.Error())
	assert.Empty(t, project)
}

func TestImport_OnImportProjectUsecase_WhenUnableToRetrieveMetadataFromRemoteRepository_ShouldReturnError(t *testing.T) {
	prMock := projectRepositoryMock{
		err: repository.ErrNoResults,
	}
	rprMock := remoteProjectRepositoryMock{
		err: repository.ErrUnexpected,
	}
	uc := usecase.NewImportProjectUsecase(prMock, rprMock)

	project, err := uc.Import(context.TODO(), "https://github.com/test/mytest")

	assert.EqualError(t, err, usecase.ErrUnableToRetrieveMetadata.Error())
	assert.Empty(t, project)
}

func TestImport_OnImportProjectUsecase_WhenUnableToCloneSourceCode_ShouldReturnError(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

// mocks
type projectRepositoryMock struct {
	project repository.Project
	err     error
}

func (m projectRepositoryMock) Add(ctx context.Context, p repository.Project) error {
	return nil
}

func (m projectRepositoryMock) GetByURL(ctx context.Context, url string) (repository.Project, error) {
	return m.project, m.err
}

type remoteProjectRepositoryMock struct {
	metadata repository.Metadata
	err      error
}

func (m remoteProjectRepositoryMock) RetrieveMetadata(ctx context.Context, url string) (repository.Metadata, error) {
	return m.metadata, m.err
}

// end mocks
