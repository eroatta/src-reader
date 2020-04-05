package splitter_test

import (
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/splitter"
	"github.com/stretchr/testify/assert"
)

func TestNewSplitterFactory_ShouldReturnSplitterAbstractFactory(t *testing.T) {
	af := splitter.NewSplitterFactory()

	assert.NotNil(t, af)
	assert.Implements(t, (*entity.SplitterAbstractFactory)(nil), af)
}

func TestGet_OnSplitterFactory_WithNotExistingAlgorithm_ShouldReturnError(t *testing.T) {
	af := splitter.NewSplitterFactory()
	got, err := af.Get("non-existing")

	assert.Nil(t, got)
	assert.Error(t, err)
}

func TestGet_OnSplitterFactory_WithConserv_ShouldReturnConservFactory(t *testing.T) {
	af := splitter.NewSplitterFactory()
	got, err := af.Get("conserv")

	assert.Implements(t, (*entity.SplitterFactory)(nil), got)
	assert.NoError(t, err)
}
