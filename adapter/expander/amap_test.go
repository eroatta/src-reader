package expander_test

import (
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/adapter/expander"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewAMAPFactory_ShouldReturnExpanderFactory(t *testing.T) {
	factory := expander.NewAMAPFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnAMAPFactory_WhenMissingScopedDeclarations_ShouldReturnError(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := expander.NewAMAPFactory()
	expander, err := factory.Make(miningResults)

	assert.Nil(t, expander)
	assert.Error(t, err)
}

func TestMake_OnAMAPFactory_WhenMissingReferenceText_ShouldReturnError(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerScopedDeclarations: miner.NewScope("test"),
	}

	factory := expander.NewAMAPFactory()
	expander, err := factory.Make(miningResults)

	assert.Nil(t, expander)
	assert.Error(t, err)
}

func TestApplicableOn_OnAMAP_ShouldReturnSamurai(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerScopedDeclarations: miner.NewScope("test"),
		entity.MinerComments:           miner.NewComments(),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(miningResults)

	assert.NotNil(t, amap)
	assert.Equal(t, "samurai", amap.ApplicableOn())
}

func TestExpand_OnAMAPWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerScopedDeclarations: miner.NewScope("test"),
		entity.MinerComments:           miner.NewComments(),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		Name: "str",
		Splits: map[string][]string{
			"gentest": []string{"str"},
		},
	}

	got := amap.Expand(ident)

	assert.Equal(t, 0, len(got))
}

func TestExpand_OnAMAPWhenNoDeclFound_ShouldReturnUnexpandedResults(t *testing.T) {
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerScopedDeclarations: miner.NewScope("test"),
		entity.MinerComments:           miner.NewComments(),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(miningResults)
	ident := entity.Identifier{
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
	miningResults := map[entity.MinerType]entity.Miner{
		entity.MinerScopedDeclarations: miner.Scope{
			Scopes: map[string]entity.ScopedDecl{
				"sb": entity.ScopedDecl{
					ID:       "sb",
					DeclType: token.FUNC,
					Comments: []string{"string buffer"},
				},
			},
		},
		entity.MinerComments: miner.NewComments(),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		Name: "sb",
		Splits: map[string][]string{
			"samurai": []string{"sb"},
		},
	}

	got := amap.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []string{"string buffer"}, got)
}
