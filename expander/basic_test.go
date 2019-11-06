package expander_test

import (
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/expander"
	"github.com/eroatta/src-reader/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewBasic_ShouldReturnBasicExpander(t *testing.T) {
	basic := expander.NewBasic(map[string]miner.Decl{})

	assert.NotNil(t, basic)
	assert.Equal(t, "basic", basic.Name())
}

func TestApplicableOn_OnBasicExpander_ShouldReturnGreedy(t *testing.T) {
	basic := expander.NewBasic(map[string]miner.Decl{})

	assert.NotNil(t, basic)
	assert.Equal(t, "greedy", basic.ApplicableOn())
}

func TestExpand_OnBasicWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	// TODO: miner.Decl should be moved to code.Decl
	var decls map[string]miner.Decl
	basic := expander.NewBasic(decls)

	ident := code.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"gentest": []string{"str"},
		},
	}

	got := basic.Expand(ident)

	assert.Equal(t, 0, len(got))
}

func TestExpand_OnBasicWhenNoDeclFound_ShouldReturnUnexpandedResults(t *testing.T) {
	var decls map[string]miner.Decl
	basic := expander.NewBasic(decls)

	ident := code.Identifier{
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
	strbuffDecl := miner.Decl{
		ID:       "strbuff",
		DeclType: token.FUNC,
		Words: map[string]struct{}{
			"string": struct{}{},
			"buffer": struct{}{},
		},
		Phrases: map[string]struct{}{},
	}
	basic := expander.NewBasic(map[string]miner.Decl{
		"strbuff": strbuffDecl,
	})

	ident := code.Identifier{
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
	strbuffDecl := miner.Decl{
		ID:       "sb",
		DeclType: token.FUNC,
		Words:    map[string]struct{}{},
		Phrases: map[string]struct{}{
			"string buffer": struct{}{},
		},
	}
	basic := expander.NewBasic(map[string]miner.Decl{
		"sb": strbuffDecl,
	})

	ident := code.Identifier{
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
