package usecase

import (
	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
)

// TODO: remove
// Expander interface is used to define a custom expander.
/*type Expander interface {
	// Name returns the name of the custom expander.
	Name() string
	// ApplicableOn defines the name of splits used as input.
	ApplicableOn() string
	// Expand performs the expansion on the token as a whole.
	Expand(ident code.Identifier) []string
}*/

// expand returns a channel of code.Identifier where each element has been processed by
// every provided Expander.
func expand(identc <-chan code.Identifier, expanders ...entity.Expander) chan code.Identifier {
	expandedc := make(chan code.Identifier)
	go func() {
		for ident := range identc {
			for _, expander := range expanders {
				if _, processable := ident.Splits[expander.ApplicableOn()]; !processable {
					continue
				}

				ident.Expansions[expander.Name()] = expander.Expand(ident)
			}

			expandedc <- ident
		}

		close(expandedc)
	}()

	return expandedc
}
