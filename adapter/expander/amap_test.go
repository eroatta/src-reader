package expander_test

import (
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/adapter/expander"
	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewAMAPFactory_ShouldReturnExpanderFactory(t *testing.T) {
	factory := expander.NewAMAPFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnAMAPFactory_WhenMissingScopedDeclarations_ShouldReturnError(t *testing.T) {
	staticInputs := map[string]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewAMAPFactory()
	expander, err := factory.Make(staticInputs, miningResults)

	assert.Nil(t, expander)
	assert.Error(t, err)
}

func TestMake_OnAMAPFactory_WhenMissingReferenceText_ShouldReturnError(t *testing.T) {
	staticInputs := map[string]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.ScopedDeclarations: miner.NewScope("test"),
	}

	factory := expander.NewAMAPFactory()
	expander, err := factory.Make(staticInputs, miningResults)

	assert.Nil(t, expander)
	assert.Error(t, err)
}

func TestApplicableOn_OnAMAP_ShouldReturnSamurai(t *testing.T) {
	staticInputs := map[string]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.ScopedDeclarations: miner.NewScope("test"),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(staticInputs, miningResults)

	assert.NotNil(t, amap)
	assert.Equal(t, "samurai", amap.ApplicableOn())
}

func TestExpand_OnAMAPWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	staticInputs := map[string]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.ScopedDeclarations: miner.NewScope("test"),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(staticInputs, miningResults)

	// TODO: miner.ScopedDecl should be moved to code.ScopedDecl
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
	staticInputs := map[string]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.ScopedDeclarations: miner.NewScope("test"),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(staticInputs, miningResults)
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
	staticInputs := map[string]interface{}{}

	miningResults := map[entity.MinerType]entity.Miner{
		entity.ScopedDeclarations: miner.Scope{
			Scopes: map[string]miner.ScopedDecl{
				"sb": miner.ScopedDecl{
					ID:       "sb",
					DeclType: token.FUNC,
					Comments: []string{"string buffer"},
				},
			},
		},
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(staticInputs, miningResults)

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
