package expander_test

import (
	"testing"

	"github.com/eroatta/src-reader/adapter/expander"
	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewNoExpansionFactory_ShouldReturnExpanderFactory(t *testing.T) {
	factory := expander.NewNoExpansionFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnNoExpansionFactory_ShouldReturnExpander(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, err := factory.Make(miningResults)

	assert.NotNil(t, noexp)
	assert.NoError(t, err)
}

func TestApplicableOn_OnNoExpansion_ShouldReturnConserv(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, _ := factory.Make(miningResults)

	assert.NotNil(t, noexp)
	assert.Equal(t, "conserv", noexp.ApplicableOn())
}

func TestExpand_OnNoexpWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		Name: "str",
		Splits: map[string]string{
			"gentest": "str",
		},
	}

	got := noexp.Expand(ident)

	assert.Equal(t, 0, len(got))
}

func TestExpand_OnNoexp_ShouldReturnSameSplittedValues(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		Name: "str",
		Splits: map[string]string{
			"conserv": "str",
		},
	}

	got := noexp.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []string{"str"}, got)
}
