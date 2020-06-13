package step

import (
	"github.com/eroatta/src-reader/entity"
)

// Split returns a channel of entity.Identifier where each element has been processed by
// every provided Splitter.
func Split(identc <-chan entity.Identifier, splitters ...entity.Splitter) chan entity.Identifier {
	splittedc := make(chan entity.Identifier)
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
