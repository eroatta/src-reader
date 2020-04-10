package expander

import (
	"fmt"
	"log"
	"strings"

	"github.com/eroatta/src-reader/adapter/miner"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/basic"
	"github.com/eroatta/token/expansion"
)

// NewBasicFactory creates a new Basic expanders factory.
func NewBasicFactory() entity.ExpanderFactory {
	return basicFactory{}
}

type basicFactory struct{}

func (f basicFactory) Make(miningResults map[entity.MinerType]entity.Miner) (entity.Expander, error) {
	declarationsMiner, ok := miningResults[entity.MinerDeclarations]
	if !ok {
		return nil, fmt.Errorf("unable to retrieve input from %s", entity.MinerDeclarations)
	}
	declarations := declarationsMiner.(*miner.Declaration).Declarations()

	return &basicExpander{
		expander:     expander{"amap"},
		declarations: declarations,
	}, nil
}

type basicExpander struct {
	expander
	declarations map[string]entity.Decl
}

// Expand receives a entity.Identifier and processes the available splits that
// can be expanded with the current algorithm.
// On Basic, we rely on the related declaration information for the identifier.
// If no declaration information can be found, we avoid trying to expand the identifier
// because results can be broad.
// If a declaration is found but several expansions are found, we handle a subset of them.
func (b basicExpander) Expand(ident entity.Identifier) []string {
	split, ok := ident.Splits[b.ApplicableOn()]
	if !ok {
		return []string{}
	}

	declarationID := ident.ID
	if ident.IsLocal() {
		declarationID = ident.Parent
	}
	decl, ok := b.declarations[declarationID]
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
		expansions := basic.Expand(token, words, phrases, basic.DefaultExpansions)
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
