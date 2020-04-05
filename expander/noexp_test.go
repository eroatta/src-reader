package expander_test

import (
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/expander"
	"github.com/stretchr/testify/assert"
)

func TestNewNoExpansionFactory_ShouldReturnExpanderFactory(t *testing.T) {
	factory := expander.NewNoExpansionFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnNoExpansionFactory_ShouldReturnExpander(t *testing.T) {
	staticInputs := map[string]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, err := factory.Make(staticInputs, miningResults)

	assert.NotNil(t, noexp)
	assert.NoError(t, err)
}

func TestApplicableOn_OnNoExpansion_ShouldReturnConserv(t *testing.T) {
	staticInputs := map[string]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, _ := factory.Make(staticInputs, miningResults)

	assert.NotNil(t, noexp)
	assert.Equal(t, "conserv", noexp.ApplicableOn())
}

func TestExpand_OnNoexpWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	staticInputs := map[string]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, _ := factory.Make(staticInputs, miningResults)

	ident := code.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"gentest": []string{"str"},
		},
	}

	got := noexp.Expand(ident)

	assert.Equal(t, 0, len(got))
}

func TestExpand_OnNoexp_ShouldReturnSameSplittedValues(t *testing.T) {
	staticInputs := map[string]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewNoExpansionFactory()
	noexp, _ := factory.Make(staticInputs, miningResults)

	ident := code.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"conserv": []string{"str"},
		},
	}

	got := noexp.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []string{"str"}, got)
}
