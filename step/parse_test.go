package step_test

import (
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/step"
	"github.com/stretchr/testify/assert"
)

func TestMerge_OnClosedChannel_ShouldReturnEmptyArray(t *testing.T) {
	parsedc := make(chan code.File)
	close(parsedc)

	got := step.Merge(parsedc)

	assert.Empty(t, got)
}

func TestMerge_OnTwoFiles_ShouldReturnTwoFiles(t *testing.T) {
	parsedc := make(chan code.File)
	go func() {
		parsedc <- code.File{}
		parsedc <- code.File{}
		close(parsedc)
	}()

	got := step.Merge(parsedc)

	assert.Equal(t, 2, len(got))
}
