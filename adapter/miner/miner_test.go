package miner_test

import (
	"testing"

	"github.com/eroatta/src-reader/adapter/miner"
	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewMinerFactory_ShouldReturnMinerAbstractFactory(t *testing.T) {
	af := miner.NewMinerFactory()

	assert.NotNil(t, af)
	assert.Implements(t, (*entity.MinerAbstractFactory)(nil), af)
}

func TestGet_OnMinerFactory_WithNotExistingAlgorithm_ShouldReturnError(t *testing.T) {
	af := miner.NewMinerFactory()
	got, err := af.Get("non-existing")

	assert.Nil(t, got)
	assert.Error(t, err)
}

func TestGet_OnMinerFactory_WithComments_ShouldReturnCommentsFactory(t *testing.T) {
	af := miner.NewMinerFactory()
	got, err := af.Get("comments")

	assert.Implements(t, (*entity.MinerFactory)(nil), got)
	assert.NoError(t, err)
}
