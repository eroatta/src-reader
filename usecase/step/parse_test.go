package step_test

import (
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/usecase/step"
	"github.com/stretchr/testify/assert"
)

func TestParse_OnClosedChannel_ShouldSendNoElements(t *testing.T) {
	filesc := make(chan entity.File)
	close(filesc)

	parsedc := step.Parse(filesc)

	var parsedFiles int
	for range parsedc {
		parsedFiles++
	}

	assert.Equal(t, 0, parsedFiles)
}

func TestParse_OnFileWithError_ShouldSendFileWithErrorMessage(t *testing.T) {
	filesc := make(chan entity.File)
	go func() {
		filesc <- entity.File{
			Name: "failing.go",
			Raw:  []byte("packaaage failing"),
		}
		close(filesc)
	}()

	parsedc := step.Parse(filesc)

	files := make([]entity.File, 0)
	for file := range parsedc {
		files = append(files, file)
	}

	assert.Equal(t, 1, len(files))
	assert.EqualError(t, files[0].Error, "failing.go:1:1: expected 'package', found packaaage")
}

func TestParse_OnTwoFiles_ShouldSendTwoParsedFilesWithSameFileset(t *testing.T) {
	filesc := make(chan entity.File)
	go func() {
		filesc <- entity.File{
			Name: "main.go",
			Raw:  []byte("package main"),
		}

		filesc <- entity.File{
			Name: "test.go",
			Raw:  []byte("package test"),
		}

		close(filesc)
	}()

	parsedc := step.Parse(filesc)

	files := make(map[string]entity.File)
	for file := range parsedc {
		files[file.Name] = file
	}

	assert.Equal(t, 2, len(files))

	mainF := files["main.go"]
	assert.NotNil(t, mainF.AST)
	assert.NotNil(t, mainF.FileSet)
	assert.NoError(t, mainF.Error)

	testF := files["test.go"]
	assert.NotNil(t, testF.AST)
	assert.NotNil(t, testF.FileSet)
	assert.NoError(t, testF.Error)

	assert.Equal(t, mainF.FileSet, testF.FileSet)
}

func TestMerge_OnClosedChannel_ShouldReturnEmptyArray(t *testing.T) {
	parsedc := make(chan entity.File)
	close(parsedc)

	got := step.Merge(parsedc)

	assert.Empty(t, got)
}

func TestMerge_OnTwoFiles_ShouldReturnTwoFiles(t *testing.T) {
	parsedc := make(chan entity.File)
	go func() {
		parsedc <- entity.File{}
		parsedc <- entity.File{}
		close(parsedc)
	}()

	got := step.Merge(parsedc)

	assert.Equal(t, 2, len(got))
}
