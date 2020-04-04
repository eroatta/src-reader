package usecase

import (
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestExpand_OnClosedChannel_ShouldSendNoElements(t *testing.T) {
	identc := make(chan code.Identifier)
	close(identc)

	expandedc := expand(identc, expander{})

	var identifiers int
	for range expandedc {
		identifiers++
	}

	assert.Equal(t, 0, identifiers)
}

func TestExpand_OnEmptyExpander_ShouldSendElementsWithoutExpansions(t *testing.T) {
	identc := make(chan code.Identifier)
	go func() {
		identc <- code.Identifier{
			Name: "crtfile",
			Splits: map[string][]string{
				"test": []string{"crt", "file"},
			},
			Expansions: make(map[string][]string),
		}
		close(identc)
	}()

	expandedc := expand(identc, []entity.Expander{}...)

	expanded := make([]code.Identifier, 0)
	for ident := range expandedc {
		expanded = append(expanded, ident)
	}

	assert.Equal(t, 1, len(expanded))
	assert.Equal(t, 0, len(expanded[0].Expansions))
}

func TestExpand_OnOneIdentifierAndTwoExpanders_ShouldSendElementsWithOneExpansion(t *testing.T) {
	identc := make(chan code.Identifier)
	go func() {
		identc <- code.Identifier{
			Name: "ctrldel",
			Splits: map[string][]string{
				"custom": []string{"ctrl", "del"},
			},
			Expansions: make(map[string][]string),
		}
		close(identc)
	}()

	custom := expander{
		name:    "custom",
		worksOn: "custom",
		efunc: func(token []string) []string {
			return []string{"control", "delete"}
		},
	}

	skipped := expander{
		name:    "skipped",
		worksOn: "none",
		efunc: func(token []string) []string {
			return []string{"should", "not", "be", "called"}
		},
	}

	expandec := expand(identc, custom, skipped)

	expandidents := make([]code.Identifier, 0)
	for ident := range expandec {
		expandidents = append(expandidents, ident)
	}

	assert.Equal(t, 1, len(expandidents))

	expansions := expandidents[0].Expansions
	assert.Equal(t, []string{"control", "delete"}, expansions["custom"])

	_, found := expansions["skipped"]
	assert.False(t, found)
}

type expander struct {
	name    string
	worksOn string
	efunc   func([]string) []string
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

func (e expander) Expand(ident code.Identifier) []string {
	if e.efunc != nil {
		return e.efunc(ident.Splits[e.worksOn])
	}

	return []string{}
}
