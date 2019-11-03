package expander_test

import (
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

// TODO: rename test
func TestExpandIdent_OnBasic_ShouldReturnAnArrayOfStrings(t *testing.T) {
	var decls map[string]miner.Decl
	basic := expander.NewBasic(decls)

	ident := code.Identifier{
		Name:   "a",
		Splits: map[string][]string{},
	}
	// TODO: change!
	got := basic.Expand(ident.Splits["greedy"])

	assert.Equal(t, 0, len(got))
	assert.EqualValues(t, []string{}, got)
}

// test expand ident on basic when no splits applicable to the expander
// test expand ident on basic when no decl for the given identifier
// test expand ident on basic for the regular scenario
