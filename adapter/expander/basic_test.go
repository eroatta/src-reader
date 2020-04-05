package expander_test

import (
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/adapter/expander"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/token/expansion"
	"github.com/stretchr/testify/assert"
)

func TestNewBasicFactory_ShouldReturnBasicExpanderFactory(t *testing.T) {
	factory := expander.NewBasicFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnBasicFactory_WhenMissingVariableDeclarations_ShouldReturnError(t *testing.T) {
	staticInputs := map[entity.InputType]interface{}{
		entity.InputDefaultExpansions: expansion.NewSetBuilder().Build(),
	}
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewBasicFactory()
	basic, err := factory.Make(staticInputs, miningResults)

	assert.Nil(t, basic)
	assert.Error(t, err)
}

func TestMake_OnBasicFactory_WhenMissingDefaultExpansions_ShouldReturnError(t *testing.T) {
	staticInputs := map[entity.InputType]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: miner.NewDeclaration(nil),
	}

	factory := expander.NewBasicFactory()
	basic, err := factory.Make(staticInputs, miningResults)

	assert.Nil(t, basic)
	assert.Error(t, err)
}

func TestApplicableOn_OnBasic_ShouldReturnGreedy(t *testing.T) {
	staticInputs := map[entity.InputType]interface{}{
		entity.InputDefaultExpansions: expansion.NewSetBuilder().Build(),
	}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: miner.NewDeclaration(nil),
	}

	factory := expander.NewBasicFactory()
	basic, err := factory.Make(staticInputs, miningResults)

	assert.NotNil(t, basic)
	assert.Equal(t, "greedy", basic.ApplicableOn())
	assert.NoError(t, err)
}

func TestExpand_OnBasicWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	// TODO: miner.Decl should be moved to code.Decl
	staticInputs := map[entity.InputType]interface{}{
		entity.InputDefaultExpansions: expansion.NewSetBuilder().Build(),
	}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: miner.NewDeclaration(nil),
	}

	factory := expander.NewBasicFactory()
	basic, _ := factory.Make(staticInputs, miningResults)

	ident := entity.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"gentest": []string{"str"},
		},
	}

	got := basic.Expand(ident)

	assert.Equal(t, 0, len(got))
}

func TestExpand_OnBasicWhenNoDeclFound_ShouldReturnUnexpandedResults(t *testing.T) {
	staticInputs := map[entity.InputType]interface{}{
		entity.InputDefaultExpansions: expansion.NewSetBuilder().Build(),
	}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: miner.NewDeclaration(nil),
	}

	factory := expander.NewBasicFactory()
	basic, _ := factory.Make(staticInputs, miningResults)

	ident := entity.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"greedy": []string{"str"},
		},
	}

	got := basic.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []string{"str"}, got)
}

func TestExpand_OnBasic_ShouldReturnExpandedResultsFromWords(t *testing.T) {
	staticInputs := map[entity.InputType]interface{}{
		entity.InputDefaultExpansions: expansion.NewSetBuilder().Build(),
	}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: miner.Declaration{
			Decls: map[string]miner.Decl{
				"strbuff": miner.Decl{
					ID:       "strbuff",
					DeclType: token.FUNC,
					Words: map[string]struct{}{
						"string": struct{}{},
						"buffer": struct{}{},
					},
					Phrases: map[string]struct{}{},
				},
			},
		},
	}

	factory := expander.NewBasicFactory()
	basic, _ := factory.Make(staticInputs, miningResults)

	ident := entity.Identifier{
		Name: "strbuff",
		Splits: map[string][]string{
			"greedy": []string{"str", "buff"},
		},
	}

	got := basic.Expand(ident)

	assert.Equal(t, 2, len(got))
	assert.EqualValues(t, []string{"string", "buffer"}, got)
}

func TestExpand_OnBasic_ShouldReturnExpandedResultsFromPhrases(t *testing.T) {
	staticInputs := map[entity.InputType]interface{}{
		entity.InputDefaultExpansions: expansion.NewSetBuilder().Build(),
	}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: miner.Declaration{
			Decls: map[string]miner.Decl{
				"sb": miner.Decl{
					ID:       "sb",
					DeclType: token.FUNC,
					Words:    map[string]struct{}{},
					Phrases: map[string]struct{}{
						"string buffer": struct{}{},
					},
				},
			},
		},
	}

	factory := expander.NewBasicFactory()
	basic, _ := factory.Make(staticInputs, miningResults)

	ident := entity.Identifier{
		Name: "sb",
		Splits: map[string][]string{
			"greedy": []string{"sb"},
		},
	}

	got := basic.Expand(ident)

	assert.Equal(t, 2, len(got))
	assert.EqualValues(t, []string{"string", "buffer"}, got)
}

func TestExpand_OnBasicWhenMultipleResults_ShouldReturnClosestThreePerWord(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}
