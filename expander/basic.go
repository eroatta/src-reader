package expander

import (
	"strings"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/src-reader/step"
	"github.com/eroatta/token/basic"
	"github.com/eroatta/token/expansion"
)

type basicExpander struct {
	expander
	declarations map[string]miner.Decl
	defaultWords expansion.Set
}

func (b basicExpander) Expand(tokens []string) []string {
	var expanded []string
	return expanded
}

func (b basicExpander) ExpandIdent(ident code.Identifier) []string {
	var expanded []string

	decl := b.declarations[ident.File]

	wordsBuilder := expansion.NewSetBuilder()
	for k := range decl.Words {
		wordsBuilder.AddStrings(k)
	}
	words := wordsBuilder.Build()

	phrases := make(map[string]string)
	for phrase := range decl.Phrases {
		// TODO
		var acron string
		phrases[acron] = strings.Join(strings.Split(phrase, " "), "-")
	}

	tokens := ident.Splits[b.ApplicableOn()]
	for _, token := range tokens {
		expansions := basic.Expand(token, words, phrases, b.defaultWords)
		expanded = append(expanded, expansions...)
	}

	return expanded
}

func (b basicExpander) ApplicableOn() string {
	return "greedy"
}

// NewBasic creates a new Basic expander.
func NewBasic(declarations map[string]miner.Decl) step.Expander {
	return basicExpander{
		expander:     expander{"basic"},
		declarations: declarations,
		defaultWords: nil,
	}
}
