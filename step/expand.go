package step

import (
	"github.com/eroatta/src-reader/code"
)

// Expander interface is used to define a custom expander.
type Expander interface {
	// Name returns the name of the custom expander.
	Name() string
	// ApplicableOn defines the name of splits used as input.
	ApplicableOn() string
	// Expand performs the expansion on the token as a whole.
	Expand(splits []string) []string
}

// Expand returns a channel of code.Identifier where each element has been processed by
// every provided Expander.
func Expand(identc <-chan code.Identifier, expanders ...Expander) chan code.Identifier {
	expandedc := make(chan code.Identifier)
	go func() {
		for ident := range identc {
			for _, expander := range expanders {
				split, processable := ident.Splits[expander.ApplicableOn()]
				if !processable {
					continue
				}

				ident.Expansions[expander.Name()] = expander.Expand(split)
			}

			expandedc <- ident
		}

		close(expandedc)
	}()

	return expandedc
}
