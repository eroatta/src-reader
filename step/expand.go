package step

import (
	"log"

	"github.com/eroatta/src-reader/code"
)

// Expand TODO
func Expand(identc <-chan code.Identifier) chan code.Identifier {
	expandedc := make(chan code.Identifier)
	go func() {
		for ident := range identc {
			log.Println("Expanding...")
			// TODO: add expansions

			expandedc <- ident
		}

		close(expandedc)
	}()

	return expandedc
}
