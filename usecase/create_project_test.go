package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewCreateProjectUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewCreateProjectUsecase(nil, nil, nil)

	assert.NotNil(t, uc)
}

func TestProcess_OnCreateProjectUsecase_ShouldReturnImportResults(t *testing.T) {
	prMock := projectRepositoryMock{
		getByURLErr: repository.ErrProjectNoResults,
	}
	rprMock := metadataRepositoryMock{
		metadata: entity.Metadata{
			Fullname: "test/mytest",
			Owner:    "test",
		},
	}
	scrMock := sourceCodeRepositoryMock{
		sourceCode: entity.SourceCode{
			Hash:     "asdasda",
			Location: "/tmp/src-code-location",
			Files:    []string{"myfile.go"},
		},
		err: nil,
	}
	uc := usecase.NewCreateProjectUsecase(prMock, rprMock, scrMock)

	project, err := uc.Process(context.TODO(), "https://github.com/test/mytest")

	assert.NoError(t, err)
	assert.Equal(t, "done", project.Status)
	assert.Equal(t, "https://github.com/test/mytest", project.URL)
	assert.Equal(t, "test/mytest", project.Metadata.Fullname)
	assert.Equal(t, "asdasda", project.SourceCode.Hash)
	assert.Equal(t, 1, len(project.SourceCode.Files))
}

func TestProcess_OnCreateProjectUsecase_WhenAlreadyImportedProject_ShouldImportResults(t *testing.T) {
	prMock := projectRepositoryMock{
		project: entity.Project{
			Status: "finished",
			URL:    "https://github.com/test/mytest",
		},
		getByURLErr: nil,
	}
	uc := usecase.NewCreateProjectUsecase(prMock, nil, nil)

	project, err := uc.Process(context.TODO(), "https://github.com/test/mytest")

	assert.NoError(t, err)
	assert.NotEmpty(t, project)
	assert.Equal(t, "finished", project.Status)
	assert.Equal(t, "https://github.com/test/mytest", project.URL)
}

func TestProcess_OnCreateProjectUsecase_WhenUnableToCheckExistingProject_ShouldReturnError(t *testing.T) {
	prMock := projectRepositoryMock{
		project:     entity.Project{},
		getByURLErr: repository.ErrProjectUnexpected,
	}
	uc := usecase.NewCreateProjectUsecase(prMock, nil, nil)

	project, err := uc.Process(context.TODO(), "https://github.com/test/mytest")

	assert.EqualError(t, err, usecase.ErrUnableToReadProject.Error())
	assert.Empty(t, project)
}

func TestProcess_OnCreateProjectUsecase_WhenUnableToRetrieveMetadataFromRemoteRepository_ShouldReturnError(t *testing.T) {
	prMock := projectRepositoryMock{
		getByURLErr: repository.ErrProjectNoResults,
	}
	rprMock := metadataRepositoryMock{
		err: repository.ErrProjectUnexpected,
	}
	uc := usecase.NewCreateProjectUsecase(prMock, rprMock, nil)

	project, err := uc.Process(context.TODO(), "https://github.com/test/mytest")

	assert.EqualError(t, err, usecase.ErrUnableToRetrieveMetadata.Error())
	assert.Empty(t, project)
}

func TestProcess_OnCreateProjectUsecase_WhenUnableToCloneSourceCode_ShouldReturnError(t *testing.T) {
	prMock := projectRepositoryMock{
		getByURLErr: repository.ErrProjectNoResults,
	}
	rprMock := metadataRepositoryMock{
		metadata: entity.Metadata{
			Fullname: "test/mytest",
			Owner:    "test",
		},
	}
	scrMock := sourceCodeRepositoryMock{
		err: repository.ErrProjectUnexpected,
	}
	uc := usecase.NewCreateProjectUsecase(prMock, rprMock, scrMock)

	project, err := uc.Process(context.TODO(), "https://github.com/test/mytest")

	assert.EqualError(t, err, usecase.ErrUnableToCloneSourceCode.Error())
	assert.Empty(t, project)
}

func TestProcess_OnCreateProjectUsecase_WhenUnableToSaveImportedProject_ShouldReturnError(t *testing.T) {
	prMock := projectRepositoryMock{
		getByURLErr: repository.ErrProjectNoResults,
		addErr:      repository.ErrProjectUnexpected,
	}
	rprMock := metadataRepositoryMock{
		metadata: entity.Metadata{
			Fullname: "test/mytest",
			Owner:    "test",
		},
	}
	scrMock := sourceCodeRepositoryMock{
		err: nil,
	}
	uc := usecase.NewCreateProjectUsecase(prMock, rprMock, scrMock)

	project, err := uc.Process(context.TODO(), "https://github.com/test/mytest")

	assert.EqualError(t, err, usecase.ErrUnableToSaveProject.Error())
	assert.Empty(t, project)
}

// mocks
type projectRepositoryMock struct {
	project     entity.Project
	getByURLErr error
	addErr      error
}

func (m projectRepositoryMock) Add(ctx context.Context, p entity.Project) error {
	return m.addErr
}

func (m projectRepositoryMock) GetByURL(ctx context.Context, url string) (entity.Project, error) {
	return m.project, m.getByURLErr
}

type metadataRepositoryMock struct {
	metadata entity.Metadata
	err      error
}

func (m metadataRepositoryMock) RetrieveMetadata(ctx context.Context, url string) (entity.Metadata, error) {
	return m.metadata, m.err
}

type sourceCodeRepositoryMock struct {
	sourceCode entity.SourceCode
	err        error
}

func (m sourceCodeRepositoryMock) Clone(ctx context.Context, fullname string, url string) (entity.SourceCode, error) {
	return m.sourceCode, m.err
}

func (m sourceCodeRepositoryMock) Remove(ctx context.Context, location string) error {
	return m.err
}

func (m sourceCodeRepositoryMock) Read(ctx context.Context, location string, filename string) ([]byte, error) {
	return []byte{}, errors.New("shouldn't be called")
}

// end mocks
