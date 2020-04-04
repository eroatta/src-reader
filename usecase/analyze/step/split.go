package step

import (
	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
)

// Split returns a channel of code.Identifier where each element has been processed by
// every provided Splitter.
func Split(identc <-chan code.Identifier, splitters ...entity.Splitter) chan code.Identifier {
	splittedc := make(chan code.Identifier)
	go func() {
		for ident := range identc {
			for _, splitter := range splitters {
				ident.Splits[splitter.Name()] = splitter.Split(ident.Name)
			}

			splittedc <- ident
		}

		close(splittedc)
	}()

	return splittedc
}
