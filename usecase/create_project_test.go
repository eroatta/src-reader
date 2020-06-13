package usecase_test

import (
	"context"
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
