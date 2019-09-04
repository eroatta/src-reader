package miner_test

import (
	"testing"

	"github.com/eroatta/src-reader/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewFunction_ShouldReturnFunctionMiner(t *testing.T) {
	miner := miner.NewFunction()

	assert.NotNil(t, miner)
}

func TestGetName_OnFunction_ShouldReturnFunction(t *testing.T) {
	miner := miner.NewFunction()

	assert.Equal(t, "function", miner.Name())
}

func TestVisit_OnFunction_ShouldReturnFunctionsWordsAndPhrases(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}
