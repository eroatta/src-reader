package expander_test

import (
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/adapter/expander"
	"github.com/eroatta/src-reader/adapter/miner"
	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewBasicFactory_ShouldReturnBasicExpanderFactory(t *testing.T) {
	factory := expander.NewBasicFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnBasicFactory_WhenMissingVariableDeclarations_ShouldReturnError(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewBasicFactory()
	basic, err := factory.Make(miningResults)

	assert.Nil(t, basic)
	assert.Error(t, err)
}

func TestApplicableOn_OnBasic_ShouldReturnGreedy(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: miner.NewDeclaration(nil),
	}

	factory := expander.NewBasicFactory()
	basic, err := factory.Make(miningResults)

	assert.NotNil(t, basic)
	assert.Equal(t, "greedy", basic.ApplicableOn())
	assert.NoError(t, err)
}

func TestExpand_OnBasicWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: miner.NewDeclaration(nil),
	}

	factory := expander.NewBasicFactory()
	basic, _ := factory.Make(miningResults)

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
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: miner.NewDeclaration(nil),
	}

	factory := expander.NewBasicFactory()
	basic, _ := factory.Make(miningResults)

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
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: &miner.Declaration{
			Decls: map[string]entity.Decl{
				"strbuff": entity.Decl{
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
	basic, _ := factory.Make(miningResults)

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
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerDeclarations: &miner.Declaration{
			Decls: map[string]entity.Decl{
				"sb": entity.Decl{
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
	basic, _ := factory.Make(miningResults)

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
