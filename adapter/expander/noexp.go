package expander

import (
	"github.com/eroatta/src-reader/entity"
)

// NewNoExpansionFactory creates a new No Expansion expander factory.
func NewNoExpansionFactory() entity.ExpanderFactory {
	return noexpFactory{}
}

type noexpFactory struct{}

func (f noexpFactory) Make(map[entity.MinerType]entity.Miner) (entity.Expander, error) {
	return noexpExpander{
		expander: expander{"noexp"},
	}, nil
}

type noexpExpander struct {
	expander
}

func (e noexpExpander) Expand(ident entity.Identifier) []entity.Expansion {
	splits, ok := ident.Splits[e.ApplicableOn()]
	if !ok {
		return []entity.Expansion{}
	}

	expansions := make([]entity.Expansion, len(splits))
	for i, split := range splits {
		expansions[i] = entity.Expansion{From: split.Value, Values: []string{split.Value}}
	}
	return expansions
}

func (e noexpExpander) ApplicableOn() string {
	return "conserv"
}
