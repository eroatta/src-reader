package step_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/usecase/analyze/step"
	"github.com/stretchr/testify/assert"
)

func TestRead_OnErrorWhileRetrievingFile_ShouldReturnFileContainingError(t *testing.T) {
	scMock := sourceCodeFileReaderMock{
		err: errors.New("error reading file"),
	}

	filesc := step.Read(context.TODO(), scMock, "/tmp/", []string{"main.go"})

	assert.NotNil(t, filesc)
	for file := range filesc {
		assert.EqualError(t, file.Error, "error reading file")
	}
}

func TestClone_OnNonGolangRepository_ShouldReturnZeroFiles(t *testing.T) {
	filesc := step.Read(context.TODO(), nil, "/tmp/", []string{"README.md"})

	assert.NotNil(t, filesc)

	numberOfFiles := 0
	for range filesc {
		numberOfFiles++
	}
	assert.Equal(t, 0, numberOfFiles)
}

func TestClone_OnGolangRepository_ShouldReturnAllGolangFiles(t *testing.T) {
	scMock := sourceCodeFileReaderMock{
		files: map[string][]byte{
			"main.go": []byte("package main"),
			"test.go": []byte("package test"),
		},
	}

	filesc := step.Read(context.TODO(), scMock, "/tmp", []string{"main.go", "test.go"})

	assert.NotNil(t, filesc)

	files := make(map[string]code.File)
	for file := range filesc {
		files[file.Name] = file
	}

	assert.Equal(t, 2, len(files))
	assert.Equal(t, []byte("package main"), files["main.go"].Raw)
	assert.Equal(t, []byte("package test"), files["test.go"].Raw)
}

type sourceCodeFileReaderMock struct {
	files map[string][]byte
	err   error
}

func (m sourceCodeFileReaderMock) Clone(ctx context.Context, fullname string, cloneURL string) (entity.SourceCode, error) {
	return entity.SourceCode{}, errors.New("shouldn't be called")
}

func (m sourceCodeFileReaderMock) Remove(ctx context.Context, location string) error {
	return errors.New("shouldn't be called")
}

func (m sourceCodeFileReaderMock) Read(ctx context.Context, location string, filename string) ([]byte, error) {
	if m.err != nil {
		return []byte{}, m.err
	}

	b, ok := m.files[filename]
	if !ok {
		return []byte{}, errors.New("not found")
	}

	return b, nil
}
