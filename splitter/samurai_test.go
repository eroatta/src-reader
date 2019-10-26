package splitter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eroatta/token/samurai"

	"github.com/eroatta/src-reader/splitter"
)

func TestNewSamurai_ShouldReturnSamuraiSplitter(t *testing.T) {
	context := samurai.NewTokenContext(nil, nil)
	splitter := splitter.NewSamurai(context)

	assert.NotNil(t, splitter)
	assert.Equal(t, "samurai", splitter.Name())
}

func TestSplit_OnSamurai_ShouldReturnAnArrayOfStrings(t *testing.T) {
	context := samurai.NewTokenContext(samurai.NewFrequencyTable(), samurai.NewFrequencyTable())
	splitter := splitter.NewSamurai(context)
	got := splitter.Split("car")

	assert.ElementsMatch(t, []string{"car"}, got)
}
