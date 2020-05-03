package expander_test

import (
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/adapter/algorithm/expander"
	"github.com/eroatta/src-reader/adapter/algorithm/miner"
	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewAMAPFactory_ShouldReturnExpanderFactory(t *testing.T) {
	factory := expander.NewAMAPFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnAMAPFactory_WhenMissingScopedDeclarations_ShouldReturnError(t *testing.T) {
	miningResults := map[string]entity.Miner{}

	factory := expander.NewAMAPFactory()
	expander, err := factory.Make(miningResults)

	assert.Nil(t, expander)
	assert.Error(t, err)
}

func TestMake_OnAMAPFactory_WhenMissingReferenceText_ShouldReturnError(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"scoped-declarations": miner.NewScope(),
	}

	factory := expander.NewAMAPFactory()
	expander, err := factory.Make(miningResults)

	assert.Nil(t, expander)
	assert.Error(t, err)
}

func TestApplicablestrOn_OnAMAP_ShouldReturnSamurai(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"scoped-declarations": miner.NewScope(),
		"comments":            miner.NewComments(),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(miningResults)

	assert.NotNil(t, amap)
	assert.Equal(t, "samurai", amap.ApplicableOn())
}

func TestExpand_OnAMAPWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"scoped-declarations": miner.NewScope(),
		"comments":            miner.NewComments(),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		ID:   "filename:main.go+++pkg:main+++declType:var+++name:str",
		Name: "str",
		Splits: map[string][]entity.Split{
			"gentest": {
				{Order: 1, Value: "str"},
			},
		},
	}

	got := amap.Expand(ident)

	assert.Equal(t, 0, len(got))
}

func TestExpand_OnAMAPWhenNoDeclFound_ShouldReturnUnexpandedResults(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"scoped-declarations": miner.NewScope(),
		"comments":            miner.NewComments(),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(miningResults)
	ident := entity.Identifier{
		ID:   "filename:main.go+++pkg:main+++declType:var+++name:str",
		Name: "str",
		Splits: map[string][]entity.Split{
			"samurai": {
				{Order: 1, Value: "str"},
			},
		},
	}

	got := amap.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []entity.Expansion{{From: "str", Values: []string{"str"}}}, got)
}

func TestExpand_OnAMAP_ShouldReturnExpandedResults(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"scoped-declarations": &miner.Scope{
			Scopes: map[string]miner.ScopedDecl{
				"filename:main.go+++pkg:main+++declType:var+++name:sb": {
					ID:       "sb",
					DeclType: token.FUNC,
					Comments: []string{"string buffer"},
				},
			},
		},
		"comments": miner.NewComments(),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		ID:   "filename:main.go+++pkg:main+++declType:var+++name:sb",
		Name: "sb",
		Splits: map[string][]entity.Split{
			"samurai": {
				{Order: 1, Value: "sb"},
			},
		},
	}

	got := amap.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []entity.Expansion{{From: "sb", Values: []string{"string buffer"}}}, got)
}

func TestExpand_OnAMAP_WhileUsingLocalIdentifier_ShouldReturnExpandedResults(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"scoped-declarations": &miner.Scope{
			Scopes: map[string]miner.ScopedDecl{
				"filename:main.go+++pkg:main+++declType:var+++name:sb": {
					ID:       "sb",
					DeclType: token.FUNC,
					Comments: []string{"string buffer"},
				},
			},
		},
		"comments": miner.NewComments(),
	}

	factory := expander.NewAMAPFactory()
	amap, _ := factory.Make(miningResults)

	ident := entity.Identifier{
		ID:   "filename:main.go+++pkg:main+++declType:var+++name:sb+++local:45",
		Name: "sb",
		Splits: map[string][]entity.Split{
			"samurai": {
				{Order: 1, Value: "sb"},
			},
		},
		Parent: "filename:main.go+++pkg:main+++declType:var+++name:sb",
	}

	got := amap.Expand(ident)

	assert.Equal(t, 1, len(got))
	assert.EqualValues(t, []entity.Expansion{{From: "sb", Values: []string{"string buffer"}}}, got)
}
