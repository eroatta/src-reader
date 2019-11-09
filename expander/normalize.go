package expander

import (
	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/step"
)

type normalizeExpander struct {
	expander
}

// Expand receives a code.Identifier and processes the available splits that
// can be expanded with the current algorithm.
// On Normalize, we rely on context information and the basics of the GenTest/Normalize algorithms.
func (n normalizeExpander) Expand(ident code.Identifier) []string {
	_, ok := ident.Splits[n.ApplicableOn()]
	if !ok {
		return []string{}
	}

	return []string{ident.Name} //gentest.Expand(ident.Name, n.simCalculator)
}

func (n normalizeExpander) ApplicableOn() string {
	return "gentest"
}

func NewNormalize() step.Expander {
	return normalizeExpander{
		expander: expander{"normalize"},
		//simCalculator:      nil, // TODO add
		//context:            nil,
		//possibleExpansions: nil,
	}
}
