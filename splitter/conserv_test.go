package splitter_test

import (
	"testing"

	"github.com/eroatta/src-reader/splitter"
	"github.com/stretchr/testify/assert"
)

func TestNewConserv_ShouldReturnConservSplitter(t *testing.T) {
	splitter := splitter.NewConserv()

	assert.NotNil(t, splitter)
	assert.Equal(t, "conserv", splitter.Name())
}
