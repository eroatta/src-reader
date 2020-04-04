package expander

import (
	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/token/expansion"
	"github.com/eroatta/token/gentest"
	"github.com/eroatta/token/lists"
)

type normalizeExpander struct {
	expander
	simCalculator      gentest.SimilarityCalculator
	declarations       map[string]miner.Decl
	possibleExpansions expansion.Set
}

// Expand receives a code.Identifier and processes the available splits that
// can be expanded with the current algorithm.
// On Normalize, we rely on context information and the basics of the GenTest/Normalize algorithms.
func (n normalizeExpander) Expand(ident code.Identifier) []string {
	split, ok := ident.Splits[n.ApplicableOn()]
	if !ok {
		return []string{}
	}

	// TODO: use key
	decl, ok := n.declarations[ident.Name]
	if !ok {
		return split
	}

	contextBuilder := lists.NewBuilder()
	for k := range decl.Words {
		contextBuilder.Add(k)
	}

	return gentest.Expand(ident.Name, n.simCalculator, contextBuilder.Build(), n.possibleExpansions)
}

func (n normalizeExpander) ApplicableOn() string {
	return "gentest"
}

// NewNormalize creates a Normalize expander with the provided Similarity Calculator and the
// given map of declarations.
func NewNormalize(simCalculator gentest.SimilarityCalculator, decls map[string]miner.Decl, expansions expansion.Set) entity.Expander {
	return normalizeExpander{
		expander:           expander{"normalize"},
		simCalculator:      simCalculator,
		declarations:       decls,
		possibleExpansions: expansions,
	}
}
