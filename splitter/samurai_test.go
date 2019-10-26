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
