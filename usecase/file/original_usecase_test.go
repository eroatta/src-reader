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

func TestNewOriginalFileUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := file.NewOriginalFileUsecase(nil, nil)

	assert.NotNil(t, uc)
}

func TestProcess_OnOriginalFileUsecase_WhenProjectNotFound_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		p:   entity.Project{},
		err: repository.ErrProjectNoResults,
	}
	uc := file.NewOriginalFileUsecase(projectRepositoryMock, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, file.ErrProjectNotFound.Error())
}

func TestProcess_OnOriginalFileUsecase_WhenErrorAccessingProject_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		p:   entity.Project{},
		err: repository.ErrProjectUnexpected,
	}
	uc := file.NewOriginalFileUsecase(projectRepositoryMock, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, file.ErrUnexpected.Error())
}

func TestProcess_OnOriginalFileUsecase_WhenFileNotFound_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		p: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "/tmp",
				Files:    []string{"test.go"},
			},
		},
		err: nil,
	}
	uc := file.NewOriginalFileUsecase(projectRepositoryMock, nil)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, file.ErrFileNotFound.Error())
}

func TestProcess_OnOriginalFileUsecase_WhenErrorAccessingFile_ShouldReturnError(t *testing.T) {
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
	uc := file.NewOriginalFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Empty(t, raw)
	assert.EqualError(t, err, file.ErrUnexpected.Error())
}

func TestProcess_OnOriginalFileUsecase_WhenValidProjectAndFile_ShouldReturnRawFile(t *testing.T) {
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
		files: map[string][]byte{
			"main.go": []byte("package main"),
		},
		err: nil,
	}
	uc := file.NewOriginalFileUsecase(projectRepositoryMock, sourceCodeRepositoryMock)

	raw, err := uc.Process(context.TODO(), "eroatta/test", "main.go")

	assert.Equal(t, []byte("package main"), raw)
	assert.NoError(t, err)
}

type projectRepositoryMock struct {
	p   entity.Project
	err error
}

func (m projectRepositoryMock) Add(ctx context.Context, project entity.Project) error {
	return errors.New("shouldn't be called")
}

func (m projectRepositoryMock) GetByURL(ctx context.Context, url string) (entity.Project, error) {
	return m.p, m.err
}

type sourceCodeRepositoryMock struct {
	files map[string][]byte
	err   error
}

func (m sourceCodeRepositoryMock) Clone(ctx context.Context, fullname string, cloneURL string) (entity.SourceCode, error) {
	return entity.SourceCode{}, errors.New("shouldn't be called")
}

func (m sourceCodeRepositoryMock) Remove(ctx context.Context, location string) error {
	return errors.New("shouldn't be called")
}

func (m sourceCodeRepositoryMock) Read(ctx context.Context, location string, filename string) ([]byte, error) {
	if m.err != nil {
		return []byte{}, m.err
	}

	b, ok := m.files[filename]
	if !ok {
		return []byte{}, errors.New("not found")
	}

	return b, nil
}
