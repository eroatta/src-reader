package step_test

import (
	"errors"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/step"
	"github.com/stretchr/testify/assert"
)

func TestClone_OnErrorWhileCloning_ShouldReturnError(t *testing.T) {
	cloner := cloner{
		repoErr: errors.New("Error cloning remote repository git@github.com:test:repo"),
	}

	repo, filesc, err := step.Clone("git@github.com:test:repo", cloner)

	assert.EqualError(t, err, "Error cloning remote repository git@github.com:test:repo")
	assert.Nil(t, repo)
	assert.Nil(t, filesc)
}

func TestClone_OnErrorWhileRetrievingFilenames_ShouldReturnError(t *testing.T) {
	cloner := cloner{
		repo:     code.Repository{Name: "github.com/test/repo"},
		filesErr: errors.New("Error retriving list of file names for git@github.com:test:repo"),
	}

	repo, filesc, err := step.Clone("git@github.com:test:repo", cloner)

	assert.EqualError(t, err, "Error retriving list of file names for git@github.com:test:repo")
	assert.Nil(t, repo)
	assert.Nil(t, filesc)
}

func TestClone_OnErrorWhileRetrievingFile_ShouldReturnFileContainingError(t *testing.T) {
	cloner := cloner{
		repo:        code.Repository{Name: "github.com/test/repo"},
		files:       []string{"main.go"},
		rawFilesErr: errors.New("Error retriving file main.go for git@github.com:test:repo"),
	}

	repo, filesc, err := step.Clone("git@github.com:test:repo", cloner)

	assert.NotNil(t, repo)
	assert.NotNil(t, filesc)
	assert.Nil(t, err)

	for file := range filesc {
		assert.EqualError(t, file.Error, "Error retriving file main.go for git@github.com:test:repo")
	}
}

func TestClone_OnNonGolangRepository_ShouldReturnZeroFiles(t *testing.T) {
	cloner := cloner{
		repo:     code.Repository{Name: "github.com/test/repo"},
		files:    []string{"README.md"},
		rawFiles: map[string][]byte{},
	}

	repo, filesc, err := step.Clone("git@github.com:test:repo", cloner)

	assert.NotNil(t, repo)
	assert.NotNil(t, filesc)
	assert.Nil(t, err)

	numberOfFiles := 0
	for range filesc {
		numberOfFiles++
	}
	assert.Equal(t, 0, numberOfFiles)
}

func TestClone_OnGolangRepository_ShouldReturnAllGolangFiles(t *testing.T) {
	cloner := cloner{
		repo:  code.Repository{Name: "github.com/test/repo"},
		files: []string{"main.go", "test.go"},
		rawFiles: map[string][]byte{
			"main.go": []byte("package main"),
			"test.go": []byte("package test"),
		},
	}

	repo, filesc, err := step.Clone("git@github.com:test:repo", cloner)

	assert.NotNil(t, repo)
	assert.NotNil(t, filesc)
	assert.Nil(t, err)

	files := make(map[string]code.File)
	for file := range filesc {
		files[file.Name] = file
	}

	assert.Equal(t, 2, len(files))
	assert.Equal(t, []byte("package main"), files["main.go"].Raw)
	assert.Equal(t, []byte("package test"), files["test.go"].Raw)
}

type cloner struct {
	repo        code.Repository
	repoErr     error
	files       []string
	filesErr    error
	rawFiles    map[string][]byte
	rawFilesErr error
}

func (c cloner) Clone(url string) (code.Repository, error) {
	if c.repoErr != nil {
		return code.Repository{}, c.repoErr
	}

	return c.repo, nil
}

func (c cloner) Filenames() ([]string, error) {
	if c.filesErr != nil {
		return []string{}, c.filesErr
	}

	return c.files, nil
}

func (c cloner) File(name string) ([]byte, error) {
	if c.rawFilesErr != nil {
		return []byte{}, c.rawFilesErr
	}

	return c.rawFiles[name], nil
}
