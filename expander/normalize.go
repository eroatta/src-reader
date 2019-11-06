package expander

import (
	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/step"
)

type normalizeExpander struct {
	expander
}

func (n normalizeExpander) ApplicableOn() string {
	return "gentest"
}

func (n normalizeExpander) Expand(ident code.Identifier) []string {
	var expanded []string
	return expanded
}

func NewNormalize() step.Expander {
	return normalizeExpander{
		expander: expander{"normalize"},
	}
}
