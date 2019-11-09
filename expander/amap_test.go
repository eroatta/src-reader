package expander_test

import (
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/expander"
	"github.com/eroatta/src-reader/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewAMAP_ShouldReturnAMAP(t *testing.T) {
	amap := expander.NewAMAP(map[string]miner.ScopedDecl{}, []string{})

	assert.NotNil(t, amap)
	assert.Equal(t, "amap", amap.Name())
}

func TestApplicableOn_OnAMAP_ShouldReturnSamurai(t *testing.T) {
	amap := expander.NewAMAP(map[string]miner.ScopedDecl{}, []string{})

	assert.NotNil(t, amap)
	assert.Equal(t, "samurai", amap.ApplicableOn())
}

func TestExpand_OnAMAPWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	// TODO: miner.ScopedDecl should be moved to code.ScopedDecl
	var scopedDecls map[string]miner.ScopedDecl
	amap := expander.NewAMAP(scopedDecls, []string{})

	ident := code.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"gentest": []string{"str"},
		},
	}

	got := amap.Expand(ident)

	assert.Equal(t, 0, len(got))
}

func TestExpand_OnAMAPWhenNoDeclFound_ShouldReturnUnexpandedResults(t *testing.T) {
	var scopedDecls map[string]miner.ScopedDecl
	amap := expander.NewAMAP(scopedDecls, []string{})

	ident := code.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"samurai": []string{"str"},
		},
	}

	got := amap.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []string{"str"}, got)
}

func TestExpand_OnAMAP_ShouldReturnExpandedResults(t *testing.T) {
	strbuffDecl := miner.ScopedDecl{
		ID:       "sb",
		DeclType: token.FUNC,
		Comments: []string{"string buffer"},
	}
	amap := expander.NewAMAP(map[string]miner.ScopedDecl{
		"sb": strbuffDecl,
	}, []string{})

	ident := code.Identifier{
		Name: "sb",
		Splits: map[string][]string{
			"samurai": []string{"sb"},
		},
	}

	got := amap.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []string{"string buffer"}, got)
}
