package usecase_test

import (
	"context"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewOriginalFileUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewOriginalFileUsecase(nil, nil)

	assert.NotNil(t, uc)
}

func TestProcess_OnOriginalFileUsecase_WhenProjectNotFound_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		project: entity.Project{},
		getErr:  repository.ErrProjectNoResults,
	}
	uc := usecase.NewOriginalFileUsecase(projectRepositoryMock, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, usecase.ErrProjectNotFound.Error())
}

func TestProcess_OnOriginalFileUsecase_WhenErrorAccessingProject_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		project: entity.Project{},
		getErr:  repository.ErrProjectUnexpected,
	}
	uc := usecase.NewOriginalFileUsecase(projectRepositoryMock, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnOriginalFileUsecase_WhenFileNotFound_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		project: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"test.go"},
			},
		},
	}
	uc := usecase.NewOriginalFileUsecase(projectRepositoryMock, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, usecase.ErrFileNotFound.Error())
}

func TestProcess_OnOriginalFileUsecase_WhenErrorAccessingFile_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		project: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"main.go"},
			},
		},
	}
	sourceCodeRepositoryMock := sourceCodeRepositoryMock{
		files: make(map[string][]byte),
		err:   repository.ErrSourceCodeUnableReadFile,
	}
	uc := usecase.NewOriginalFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnOriginalFileUsecase_WhenValidProjectAndFile_ShouldReturnRawFile(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		project: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"main.go"},
			},
		},
	}
	sourceCodeRepositoryMock := sourceCodeRepositoryMock{
		files: map[string][]byte{
			"main.go": []byte("package main"),
		},
	}
	uc := usecase.NewOriginalFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Equal(t, []byte("package main"), raw)
	assert.NoError(t, err)
}
