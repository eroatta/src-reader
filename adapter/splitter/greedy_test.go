package splitter_test

import (
	"testing"

	"github.com/eroatta/src-reader/adapter/splitter"
	"github.com/stretchr/testify/assert"
)

func TestNewGreedyFactory_ShouldReturnGreedySplitterFactory(t *testing.T) {
	factory := splitter.NewGreedyFactory()

	assert.NotNil(t, factory)

}

func TestSplit_OnGreedy_ShouldReturnAnArrayOfStrings(t *testing.T) {
	factory := splitter.NewGreedyFactory()
	splitter, _ := factory.Make(nil)

	got := splitter.Split("car")

	assert.Equal(t, "greedy", splitter.Name())
	assert.ElementsMatch(t, []string{"car"}, got)
}
