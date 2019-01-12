package extractors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: use a real test... not this piece of shit
func TestProcess_OnSamuraiExtractor_ShouldNotFail(t *testing.T) {
	samurai := NewSamuraiExtractor()

	Process(samurai, nil)

	assert.True(t, true)
}
