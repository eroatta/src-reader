package usecase_test

import (
	"context"
	"testing"

	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewImportProjectUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewImportProjectUsecase(nil, nil, nil)

	assert.NotNil(t, uc)
}

func TestImport_OnImportProjectUsecase_ShouldReturnImportResults(t *testing.T) {
	prMock := projectRepositoryMock{
		getByURLErr: repository.ErrNoResults,
	}
	rprMock := remoteProjectRepositoryMock{
		metadata: repository.Metadata{},
	}
	scrMock := sourceCodeRepositoryMock{
		err: nil,
	}
	uc := usecase.NewImportProjectUsecase(prMock, rprMock, scrMock)

	project, err := uc.Import(context.TODO(), "https://github.com/test/mytest")

	assert.NoError(t, err)
	assert.Equal(t, "done", project.Status)
}

func TestImport_OnImportProjectUsecase_WhenAlreadyImportedProject_ShouldImportResults(t *testing.T) {
	prMock := projectRepositoryMock{
		project: repository.Project{
			Status: "finished",
		},
		getByURLErr: nil,
	}
	uc := usecase.NewImportProjectUsecase(prMock, nil, nil)

	project, err := uc.Import(context.TODO(), "https://github.com/test/mytest")

	assert.NoError(t, err)
	assert.NotEmpty(t, project)
	assert.Equal(t, "finished", project.Status)
}

func TestImport_OnImportProjectUsecase_WhenUnableToCheckExistingProject_ShouldReturnError(t *testing.T) {
	prMock := projectRepositoryMock{
		project:     repository.Project{},
		getByURLErr: repository.ErrUnexpected,
	}
	uc := usecase.NewImportProjectUsecase(prMock, nil, nil)

	project, err := uc.Import(context.TODO(), "https://github.com/test/mytest")

	assert.EqualError(t, err, usecase.ErrUnableToReadProject.Error())
	assert.Empty(t, project)
}

func TestImport_OnImportProjectUsecase_WhenUnableToRetrieveMetadataFromRemoteRepository_ShouldReturnError(t *testing.T) {
	prMock := projectRepositoryMock{
		getByURLErr: repository.ErrNoResults,
	}
	rprMock := remoteProjectRepositoryMock{
		err: repository.ErrUnexpected,
	}
	uc := usecase.NewImportProjectUsecase(prMock, rprMock, nil)

	project, err := uc.Import(context.TODO(), "https://github.com/test/mytest")

	assert.EqualError(t, err, usecase.ErrUnableToRetrieveMetadata.Error())
	assert.Empty(t, project)
}

func TestImport_OnImportProjectUsecase_WhenUnableToCloneSourceCode_ShouldReturnError(t *testing.T) {
	prMock := projectRepositoryMock{
		getByURLErr: repository.ErrNoResults,
	}
	rprMock := remoteProjectRepositoryMock{
		metadata: repository.Metadata{},
	}
	scrMock := sourceCodeRepositoryMock{
		err: repository.ErrUnexpected,
	}
	uc := usecase.NewImportProjectUsecase(prMock, rprMock, scrMock)

	project, err := uc.Import(context.TODO(), "https://github.com/test/mytest")

	assert.EqualError(t, err, usecase.ErrUnableToCloneSourceCode.Error())
	assert.Empty(t, project)
}

func TestImport_OnImportProjectUsecase_WhenUnableToSaveImportedProject_ShouldReturnError(t *testing.T) {
	prMock := projectRepositoryMock{
		getByURLErr: repository.ErrNoResults,
		addErr:      repository.ErrUnexpected,
	}
	rprMock := remoteProjectRepositoryMock{
		metadata: repository.Metadata{},
	}
	scrMock := sourceCodeRepositoryMock{
		err: nil,
	}
	uc := usecase.NewImportProjectUsecase(prMock, rprMock, scrMock)

	project, err := uc.Import(context.TODO(), "https://github.com/test/mytest")

	assert.EqualError(t, err, usecase.ErrUnableToSaveProject.Error())
	assert.Empty(t, project)
}

// mocks
type projectRepositoryMock struct {
	project     repository.Project
	getByURLErr error
	addErr      error
}

func (m projectRepositoryMock) Add(ctx context.Context, p repository.Project) error {
	return m.addErr
}

func (m projectRepositoryMock) GetByURL(ctx context.Context, url string) (repository.Project, error) {
	return m.project, m.getByURLErr
}

type remoteProjectRepositoryMock struct {
	metadata repository.Metadata
	err      error
}

func (m remoteProjectRepositoryMock) RetrieveMetadata(ctx context.Context, url string) (repository.Metadata, error) {
	return m.metadata, m.err
}

type sourceCodeRepositoryMock struct {
	err error
}

func (m sourceCodeRepositoryMock) Clone(ctx context.Context, url string) error {
	return m.err
}

// end mocks
