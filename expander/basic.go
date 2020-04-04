package expander

import (
	"log"
	"strings"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/token/basic"
	"github.com/eroatta/token/expansion"
)

type basicExpander struct {
	expander
	declarations map[string]miner.Decl
	defaultWords expansion.Set
}

// Expand receives a code.Identifier and processes the available splits that
// can be expanded with the current algorithm.
// On Basic, we rely on the related declaration information for the identifier.
// If no declaration information can be found, we avoid trying to expand the identifier
// because results can be broad.
// If a declaration is found but several expansions are found, we handle a subset of them.
func (b basicExpander) Expand(ident code.Identifier) []string {
	split, ok := ident.Splits[b.ApplicableOn()]
	if !ok {
		return []string{}
	}

	// TODO: use key
	decl, ok := b.declarations[ident.Name]
	if !ok {
		return split
	}

	wordsBuilder := expansion.NewSetBuilder()
	for k := range decl.Words {
		wordsBuilder.AddStrings(k)
	}
	words := wordsBuilder.Build()

	phrases := make(map[string]string)
	for phrase := range decl.Phrases {
		var acron strings.Builder
		for _, word := range strings.Split(phrase, " ") {
			acron.WriteByte(word[0])
		}

		phrases[acron.String()] = strings.ReplaceAll(phrase, " ", "-")
	}

	var expanded []string
	for _, token := range split {
		expansions := basic.Expand(token, words, phrases, b.defaultWords)
		if len(expansions) == 0 {
			expansions = []string{token}
		}

		if len(expansions) > 1 {
			// TODO: sort them
			log.Println("multiple expansions...")
		}

		expanded = append(expanded, expansions...)
	}

	return expanded
}

func (b basicExpander) ApplicableOn() string {
	return "greedy"
}

// NewBasic creates a new Basic expander.
func NewBasic(declarations map[string]miner.Decl) entity.Expander {
	return basicExpander{
		expander:     expander{"basic"},
		declarations: declarations,
		defaultWords: basic.DefaultExpansions,
	}
}
