package step

import (
	"fmt"
	"log"

	"github.com/eroatta/src-reader/code"
)

// Storer interface is used to define a custom storer.
type Storer interface {
	Save(ident code.Identifier) error
}

// Store receives a channel of code.Identifier and stores each element using the storer.
// Unsaved identifiers are returned.
func Store(identc <-chan code.Identifier, storer Storer) []code.Identifier {
	failedIdent := make([]code.Identifier, 0)
	for ident := range identc {
		err := storer.Save(ident)
		if err != nil {
			log.Println(fmt.Sprintf("Error saving identifier %s, from file %s and position %v",
				ident.Name, ident.File, ident.Position))
			ident.Error = err

			failedIdent = append(failedIdent, ident)
		}
	}

	return failedIdent
}
