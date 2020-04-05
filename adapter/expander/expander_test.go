package expander_test

import (
	"testing"

	"github.com/eroatta/src-reader/adapter/expander"
	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewExpanderFactory_ShouldReturnExpanderAbstractFactory(t *testing.T) {
	af := expander.NewExpanderFactory()

	assert.NotNil(t, af)
	assert.Implements(t, (*entity.ExpanderAbstractFactory)(nil), af)
}

func TestGet_OnExpanderFactory_WithNotExistingAlgorithm_ShouldReturnError(t *testing.T) {
	af := expander.NewExpanderFactory()
	got, err := af.Get("non-existing")

	assert.Nil(t, got)
	assert.Error(t, err)
}

func TestGet_OnExpanderFactory_WithAMAP_ShouldReturnAMAPFactory(t *testing.T) {
	af := expander.NewExpanderFactory()
	got, err := af.Get("amap")

	assert.Implements(t, (*entity.ExpanderFactory)(nil), got)
	assert.NoError(t, err)
}
