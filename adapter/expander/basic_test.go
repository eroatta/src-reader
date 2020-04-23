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
	miningResults := map[string]entity.Miner{}

	factory := expander.NewBasicFactory()
	basic, err := factory.Make(miningResults)

	assert.Nil(t, basic)
	assert.Error(t, err)
}

func TestApplicableOn_OnBasic_ShouldReturnGreedy(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"declarations": miner.NewDeclaration(nil),
	}

	factory := expander.NewBasicFactory()
	basic, err := factory.Make(miningResults)

	assert.NotNil(t, basic)
	assert.Equal(t, "greedy", basic.ApplicableOn())
	assert.NoError(t, err)
}

func TestExpand_OnBasicWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"declarations": miner.NewDeclaration(nil),
	}

	factory := expander.NewBasicFactory()
	basic, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		ID:   "filename:main.go+++pkg:main+++declType:var+++name:str",
		Name: "str",
		Splits: map[string][]entity.Split{
			"gentest": []entity.Split{
				{Order: 1, Value: "str"},
			},
		},
	}

	got := basic.Expand(ident)

	assert.Equal(t, 0, len(got))
}

func TestExpand_OnBasicWhenNoDeclFound_ShouldReturnUnexpandedResults(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"declarations": miner.NewDeclaration(nil),
	}

	factory := expander.NewBasicFactory()
	basic, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		ID:   "filename:main.go+++pkg:main+++declType:var+++name:str",
		Name: "str",
		Splits: map[string][]entity.Split{
			"greedy": []entity.Split{
				{Order: 1, Value: "str"},
			},
		},
	}

	got := basic.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []entity.Expansion{{From: "str", Values: []string{"str"}}}, got)
}

func TestExpand_OnBasic_ShouldReturnExpandedResultsFromWords(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"declarations": &miner.Declaration{
			Decls: map[string]entity.Decl{
				"filename:main.go+++pkg:main+++declType:var+++name:strbuff": entity.Decl{
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
		ID:   "filename:main.go+++pkg:main+++declType:var+++name:strbuff",
		Name: "strbuff",
		Splits: map[string][]entity.Split{
			"greedy": []entity.Split{
				{Order: 1, Value: "str"},
				{Order: 2, Value: "buff"},
			},
		},
	}

	got := basic.Expand(ident)

	assert.Equal(t, 2, len(got))
	assert.EqualValues(t, []entity.Expansion{{From: "str", Values: []string{"string"}},
		{From: "buff", Values: []string{"buffer"}}}, got)
}

func TestExpand_OnBasic_WhileUsingLocalVariables_ShouldReturnExpandedResultsFromWords(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"declarations": &miner.Declaration{
			Decls: map[string]entity.Decl{
				"filename:main.go+++pkg:main+++declType:var+++name:strbuff": entity.Decl{
					ID:       "filename:main.go+++pkg:main+++declType:var+++name:strbuff",
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
		ID:   "filename:main.go+++pkg:main+++declType:var+++name:sb+++local:43",
		Name: "sb",
		Splits: map[string][]entity.Split{
			"greedy": []entity.Split{
				{Order: 1, Value: "s"},
				{Order: 2, Value: "b"},
			},
		},
		Parent: "filename:main.go+++pkg:main+++declType:var+++name:strbuff",
	}

	got := basic.Expand(ident)

	assert.Equal(t, 2, len(got))
	assert.EqualValues(t, []entity.Expansion{{From: "s", Values: []string{"string"}},
		{From: "b", Values: []string{"buffer"}}}, got)
}

func TestExpand_OnBasic_ShouldReturnExpandedResultsFromPhrases(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"declarations": &miner.Declaration{
			Decls: map[string]entity.Decl{
				"filename:main.go+++pkg:main+++declType:var+++name:sb": entity.Decl{
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
		ID:   "filename:main.go+++pkg:main+++declType:var+++name:sb",
		Name: "sb",
		Splits: map[string][]entity.Split{
			"greedy": []entity.Split{
				{Order: 1, Value: "sb"},
			},
		},
	}

	got := basic.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []entity.Expansion{{From: "sb", Values: []string{"string buffer"}}}, got)
}

func TestExpand_OnBasicWhenMultipleResults_ShouldReturnClosestThreePerWord(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"declarations": &miner.Declaration{
			Decls: map[string]entity.Decl{
				"filename:main.go+++pkg:main+++declType:var+++name:contrl": entity.Decl{
					ID:       "contrl",
					DeclType: token.FUNC,
					Words:    map[string]struct{}{},
					Phrases:  map[string]struct{}{},
				},
			},
		},
	}

	factory := expander.NewBasicFactory()
	basic, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		ID:   "filename:main.go+++pkg:main+++declType:var+++name:contrl",
		Name: "contrl",
		Splits: map[string][]entity.Split{
			"greedy": []entity.Split{
				{Order: 1, Value: "contrl"},
			},
		},
	}

	got := basic.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []entity.Expansion{{From: "contrl", Values: []string{"control", "control", "contrail"}}}, got)
}
