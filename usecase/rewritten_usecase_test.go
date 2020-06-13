package usecase_test

import (
	"context"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewRewrittenFileUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewRewrittenFileUsecase(nil, nil, nil)

	assert.NotNil(t, uc)
}

func TestProcess_OnRewrittenFileUsecase_WhenProjectNotFound_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		project:     entity.Project{},
		getByURLErr: repository.ErrProjectNoResults,
	}
	uc := usecase.NewRewrittenFileUsecase(projectRepositoryMock, nil, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, usecase.ErrProjectNotFound.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenErrorAccessingProject_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		project:     entity.Project{},
		getByURLErr: repository.ErrProjectUnexpected,
	}
	uc := usecase.NewRewrittenFileUsecase(projectRepositoryMock, nil, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenFileNotFound_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		project: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"test.go"},
			},
		},
	}
	uc := usecase.NewRewrittenFileUsecase(projectRepositoryMock, nil, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, usecase.ErrFileNotFound.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenErrorAccessingFile_ShouldReturnError(t *testing.T) {
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
	uc := usecase.NewRewrittenFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenNoIdentifiers_ShouldReturnError(t *testing.T) {
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
	uc := usecase.NewRewrittenFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenErrorAccessingIdentifiers_ShouldReturnError(t *testing.T) {
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
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{},
		err:    repository.ErrIdentifierUnexpected,
	}
	uc := usecase.NewRewrittenFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock, identifierRepositoryMock)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenValidProjectAndFile_ShouldReturnRawFile(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		project: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"main.go"},
			},
		},
	}
	content := []byte(`
		package main

		func main() {}
	`)
	sourceCodeRepositoryMock := sourceCodeRepositoryMock{
		files: map[string][]byte{
			"main.go": content,
		},
		err: nil,
	}
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{
			{
				Name: "main",
				Normalization: entity.Normalization{
					Word: "changed",
				},
			},
		},
		err: nil,
	}
	uc := usecase.NewRewrittenFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock, identifierRepositoryMock)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	expected := []byte(`package main

func changed() {}
`)
	assert.Equal(t, string(expected), string(raw))
	assert.NoError(t, err)
}
