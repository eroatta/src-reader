package step

import (
	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
)

// Expand returns a channel of code.Identifier where each element has been processed by
// every provided Expander.
func Expand(identc <-chan code.Identifier, expanders ...entity.Expander) chan code.Identifier {
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
