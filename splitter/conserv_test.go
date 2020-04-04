package splitter_test

import (
	"testing"

	"github.com/eroatta/src-reader/splitter"
	"github.com/stretchr/testify/assert"
)

func TestNewConservFactory_ShouldReturnConservSplitterFactory(t *testing.T) {
	factory := splitter.NewConservFactory()

	assert.NotNil(t, factory)
}

func TestSplit_OnConserv_ShouldReturnAnArrayOfStrings(t *testing.T) {
	factory := splitter.NewConservFactory()
	splitter, _ := factory.Make(nil, nil)

	got := splitter.Split("car")

	assert.Equal(t, "conserv", splitter.Name())
	assert.ElementsMatch(t, []string{"car"}, got)
}
