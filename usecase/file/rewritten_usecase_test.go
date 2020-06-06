package file_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase/file"
	"github.com/stretchr/testify/assert"
)

func TestNewRewrittenFileUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := file.NewRewrittenFileUsecase(nil, nil, nil)

	assert.NotNil(t, uc)
}

func TestProcess_OnRewrittenFileUsecase_WhenProjectNotFound_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		p:   entity.Project{},
		err: repository.ErrProjectNoResults,
	}
	uc := file.NewRewrittenFileUsecase(projectRepositoryMock, nil, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, file.ErrProjectNotFound.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenErrorAccessingProject_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		p:   entity.Project{},
		err: repository.ErrProjectUnexpected,
	}
	uc := file.NewRewrittenFileUsecase(projectRepositoryMock, nil, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, file.ErrUnexpected.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenFileNotFound_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		p: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"test.go"},
			},
		},
		err: nil,
	}
	uc := file.NewRewrittenFileUsecase(projectRepositoryMock, nil, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, file.ErrFileNotFound.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenErrorAccessingFile_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		p: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"main.go"},
			},
		},
		err: nil,
	}
	sourceCodeRepositoryMock := sourceCodeRepositoryMock{
		files: make(map[string][]byte),
		err:   repository.ErrSourceCodeUnableReadFile,
	}
	uc := file.NewRewrittenFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, file.ErrUnexpected.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenNoIdentifiers_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		p: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"main.go"},
			},
		},
		err: nil,
	}
	sourceCodeRepositoryMock := sourceCodeRepositoryMock{
		files: make(map[string][]byte),
		err:   repository.ErrSourceCodeUnableReadFile,
	}
	uc := file.NewRewrittenFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, file.ErrIdentifiersNotFound.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenErrorAccessingIdentifiers_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		p: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"main.go"},
			},
		},
		err: nil,
	}
	sourceCodeRepositoryMock := sourceCodeRepositoryMock{
		files: make(map[string][]byte),
		err:   repository.ErrSourceCodeUnableReadFile,
	}
	identifierRepositoryMock := identifierRepositoryMock{
		idents: []entity.Identifier{},
		err:    repository.ErrIdentifierUnexpected,
	}
	uc := file.NewRewrittenFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock, identifierRepositoryMock)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, file.ErrUnexpected.Error())
}

func TestProcess_OnRewrittenFileUsecase_WhenValidProjectAndFile_ShouldReturnRawFile(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		p: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"main.go"},
			},
		},
		err: nil,
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
	uc := file.NewRewrittenFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock, identifierRepositoryMock)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	expected := []byte(`package main

func changed() {}
`)
	assert.Equal(t, string(expected), string(raw))
	assert.NoError(t, err)
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

func (i identifierRepositoryMock) FindAllByProjectAndFile(ctx context.Context, projectRef string, filename string) ([]entity.Identifier, error) {
	return i.idents, i.err
}
