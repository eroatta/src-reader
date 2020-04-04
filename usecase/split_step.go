package usecase

import (
	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
)

// TODO: remove
// Splitter interface is used to define a custom splitter.
/*type Splitter interface {
	// Name returns the name of the custom splitter.
	Name() string
	// Split returns the split identifier.
	Split(token string) []string //TODO: we should return a string, divided by something like hyphen
}*/

// Split returns a channel of code.Identifier where each element has been processed by
// every provided Splitter.
func split(identc <-chan code.Identifier, splitters ...entity.Splitter) chan code.Identifier {
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
