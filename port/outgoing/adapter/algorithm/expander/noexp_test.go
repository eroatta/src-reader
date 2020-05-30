package expander_test

import (
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/port/outgoing/adapter/algorithm/expander"
	"github.com/stretchr/testify/assert"
)

func TestNewNoExpansionFactory_ShouldReturnExpanderFactory(t *testing.T) {
	factory := expander.NewNoExpansionFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnNoExpansionFactory_ShouldReturnExpander(t *testing.T) {
	miningResults := map[string]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, err := factory.Make(miningResults)

	assert.NotNil(t, noexp)
	assert.NoError(t, err)
}

func TestApplicableOn_OnNoExpansion_ShouldReturnConserv(t *testing.T) {
	miningResults := map[string]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, _ := factory.Make(miningResults)

	assert.NotNil(t, noexp)
	assert.Equal(t, "conserv", noexp.ApplicableOn())
}

func TestExpand_OnNoexpWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	miningResults := map[string]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		Name: "str",
		Splits: map[string][]entity.Split{
			"gentest": {
				{Order: 1, Value: "str"},
			},
		},
	}

	got := noexp.Expand(ident)

	assert.Equal(t, 0, len(got))
}

func TestExpand_OnNoexp_ShouldReturnSameSplittedValues(t *testing.T) {
	miningResults := map[string]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, _ := factory.Make(miningResults)

	split := entity.Split{Order: 1, Value: "str"}
	ident := entity.Identifier{
		Name: "str",
		Splits: map[string][]entity.Split{
			"conserv": {split},
		},
	}

	got := noexp.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []entity.Expansion{{Order: 1, SplittingAlgorithm: "conserv", From: "str", Values: []string{"str"}}}, got)
}
