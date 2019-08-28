package step

import (
	"github.com/eroatta/src-reader/code"
)

// Storer TODO
type Storer interface {
	Save(ident code.Identifier) error
}

// Store receives a channel of Identifier and stores each element on the given storer.
func Store(identc <-chan code.Identifier, str Storer) []code.Identifier {
	failedIdent := make([]code.Identifier, 0)
	for i := range identc {
		err := str.Save(i)
		if err != nil {
			// log.Println(err)
			i.Error = err
			failedIdent = append(failedIdent, i)
		}
	}

	return failedIdent
}
