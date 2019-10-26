package splitter_test

import (
	"testing"

	"github.com/eroatta/src-reader/splitter"
	"github.com/stretchr/testify/assert"
)

func TestNewGreedy_ShouldReturnGreedySplitter(t *testing.T) {
	splitter := splitter.NewGreedy()

	assert.NotNil(t, splitter)
	assert.Equal(t, "greedy", splitter.Name())
}

func TestSplit_OnGreedy_ShouldReturnAnArrayOfStrings(t *testing.T) {
	splitter := splitter.NewGreedy()
	got := splitter.Split("car")

	assert.ElementsMatch(t, []string{"car"}, got)
}
