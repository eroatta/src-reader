package splitter_test

import (
	"testing"

	"github.com/eroatta/src-reader/adapter/algorithm/splitter"
	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewConservFactory_ShouldReturnConservSplitterFactory(t *testing.T) {
	factory := splitter.NewConservFactory()

	assert.NotNil(t, factory)
}

func TestSplit_OnConserv_ShouldReturnSplit(t *testing.T) {
	factory := splitter.NewConservFactory()
	splitter, _ := factory.Make(nil)

	got := splitter.Split("car")

	assert.Equal(t, "conserv", splitter.Name())
	assert.Equal(t, []entity.Split{{Order: 1, Value: "car"}}, got)
}
