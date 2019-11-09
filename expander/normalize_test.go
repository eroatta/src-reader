package expander_test

import (
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/expander"
	"github.com/stretchr/testify/assert"
)

func TestNewNormalize_ShouldReturnNormalize(t *testing.T) {
	normalize := expander.NewNormalize()

	assert.NotNil(t, normalize)
	assert.Equal(t, "normalize", normalize.Name())
}

func TestApplicableOn_OnNormalize_ShouldReturnGenTest(t *testing.T) {
	normalize := expander.NewNormalize()

	assert.NotNil(t, normalize)
	assert.Equal(t, "gentest", normalize.ApplicableOn())
}

func TestExpand_OnNormalizeWhenNoSplitsApplicable_ShouldReturnEmptyResults(t *testing.T) {
	normalize := expander.NewNormalize()

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

func TestExpand_OnNormalize_ShouldReturnExpansions(t *testing.T) {
	normalize := expander.NewNormalize()

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
