package step_test

import (
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/usecase/analyze/step"
	"github.com/stretchr/testify/assert"
)

func TestExpand_OnClosedChannel_ShouldSendNoElements(t *testing.T) {
	identc := make(chan entity.Identifier)
	close(identc)

	expandedc := step.Expand(identc, expander{})

	var identifiers int
	for range expandedc {
		identifiers++
	}

	assert.Equal(t, 0, identifiers)
}

func TestExpand_OnEmptyExpander_ShouldSendElementsWithoutExpansions(t *testing.T) {
	identc := make(chan entity.Identifier)
	go func() {
		identc <- entity.Identifier{
			Name: "crtfile",
			Splits: map[string][]entity.Split{
				"test": {
					{Order: 1, Value: "crt"},
					{Order: 2, Value: "file"},
				},
			},
			Expansions: make(map[string][]entity.Expansion),
		}
		close(identc)
	}()

	expandedc := step.Expand(identc, []entity.Expander{}...)

	expanded := make([]entity.Identifier, 0)
	for ident := range expandedc {
		expanded = append(expanded, ident)
	}

	assert.Equal(t, 1, len(expanded))
	assert.Equal(t, 0, len(expanded[0].Expansions))
}

func TestExpand_OnOneIdentifierAndTwoExpanders_ShouldSendElementsWithOneExpansion(t *testing.T) {
	identc := make(chan entity.Identifier)
	go func() {
		identc <- entity.Identifier{
			Name: "ctrldel",
			Splits: map[string][]entity.Split{
				"custom": {
					{Order: 1, Value: "ctrl"},
					{Order: 2, Value: "del"},
				},
			},
			Expansions: make(map[string][]entity.Expansion),
		}
		close(identc)
	}()

	custom := expander{
		name:    "custom",
		worksOn: "custom",
		efunc: func([]entity.Split) []entity.Expansion {
			return []entity.Expansion{
				{From: "ctrl", Values: []string{"control"}},
				{From: "del", Values: []string{"delete"}},
			}
		},
	}

	skipped := expander{
		name:    "skipped",
		worksOn: "none",
		efunc: func([]entity.Split) []entity.Expansion {
			return []entity.Expansion{
				{From: "should", Values: []string{"should"}},
				{From: "not", Values: []string{"not"}},
				{From: "be", Values: []string{"be"}},
				{From: "called", Values: []string{"called"}},
			}
		},
	}

	expandec := step.Expand(identc, custom, skipped)

	expandidents := make([]entity.Identifier, 0)
	for ident := range expandec {
		expandidents = append(expandidents, ident)
	}

	assert.Equal(t, 1, len(expandidents))

	expansions := expandidents[0].Expansions
	assert.Equal(t, []entity.Expansion{{From: "ctrl", Values: []string{"control"}},
		{From: "del", Values: []string{"delete"}}}, expansions["custom"])

	_, found := expansions["skipped"]
	assert.False(t, found)
}

type expander struct {
	name    string
	worksOn string
	efunc   func([]entity.Split) []entity.Expansion
}

func (e expander) Name() string {
	if e.name != "" {
		return e.name
	}

	return "test"
}

func (e expander) ApplicableOn() string {
	return e.worksOn
}

func (e expander) Expand(ident entity.Identifier) []entity.Expansion {
	if e.efunc != nil {
		return e.efunc(ident.Splits[e.worksOn])
	}

	return []entity.Expansion{}
}
