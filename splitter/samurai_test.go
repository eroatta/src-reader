package splitter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eroatta/token/samurai"

	"github.com/eroatta/src-reader/splitter"
)

func TestNewSamurai_ShouldReturnSamuraiSplitter(t *testing.T) {
	splitter := splitter.NewSamurai(nil, nil)

	assert.NotNil(t, splitter)
	assert.Equal(t, "samurai", splitter.Name())
}

func TestSplit_OnSamurai_ShouldReturnAnArrayOfStrings(t *testing.T) {
	splitter := splitter.NewSamurai(samurai.NewFrequencyTable(), samurai.NewFrequencyTable())
	got := splitter.Split("car")

	assert.ElementsMatch(t, []string{"car"}, got)
}
