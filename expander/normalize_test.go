package expander_test

import (
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/expander"
	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/token/expansion"
	"github.com/stretchr/testify/assert"
)

func TestNewNormalize_ShouldReturnNormalize(t *testing.T) {
	var simCalculatorMock similarityCalculatorMock
	var decls map[string]miner.Decl
	var expansions expansion.Set
	normalize := expander.NewNormalize(simCalculatorMock, decls, expansions)

	assert.NotNil(t, normalize)
	assert.Equal(t, "normalize", normalize.Name())
}

func TestApplicableOn_OnNormalize_ShouldReturnGenTest(t *testing.T) {
	var simCalculatorMock similarityCalculatorMock
	var decls map[string]miner.Decl
	var expansions expansion.Set
	normalize := expander.NewNormalize(simCalculatorMock, decls, expansions)

	assert.NotNil(t, normalize)
	assert.Equal(t, "gentest", normalize.ApplicableOn())
}

func TestExpand_OnNormalizeWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	var simCalculatorMock similarityCalculatorMock
	var decls map[string]miner.Decl
	var expansions expansion.Set
	normalize := expander.NewNormalize(simCalculatorMock, decls, expansions)

	ident := code.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"greedy": []string{"str"},
		},
	}

	got := normalize.Expand(ident)

	assert.Equal(t, 0, len(got))
	assert.EqualValues(t, []string{}, got)
}

func TestExpand_OnNormalizeWhenNoDeclFound_ShouldReturnUnexpandedResults(t *testing.T) {
	var simCalculatorMock similarityCalculatorMock
	var decls map[string]miner.Decl
	var expansions expansion.Set
	normalize := expander.NewNormalize(simCalculatorMock, decls, expansions)

	ident := code.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"gentest": []string{"str"},
		},
	}

	got := normalize.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []string{"str"}, got)
}

func TestExpand_OnNormalize_ShouldReturnExpansions(t *testing.T) {
	simCalculatorMock := similarityCalculatorMock{
		"string-builder": 1.0,
	}

	decls := map[string]miner.Decl{
		"str": miner.Decl{
			ID:       "str",
			DeclType: token.FUNC,
			Words: map[string]struct{}{
				"builder": struct{}{},
			},
			Phrases: map[string]struct{}{},
		},
	}

	expansions := expansion.NewSetBuilder().AddStrings("string").Build()

	normalize := expander.NewNormalize(simCalculatorMock, decls, expansions)

	ident := code.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"gentest": []string{"str"},
		},
	}

	got := normalize.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []string{"string"}, got)
}

type similarityCalculatorMock map[string]float64

func (s similarityCalculatorMock) Similarity(word string, another string) float64 {
	var key string
	if word < another {
		key = word + "-" + another
	} else {
		key = another + "-" + word
	}

	return s[key]
}
