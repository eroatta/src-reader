package step

import (
	"log"

	"github.com/eroatta/src-reader/code"
)

// Storer TODO
type Storer interface {
	Save(code.Identifier) error
}

// Store receives a channel of Identifier and stores each element on the given storer.
func Store(identc <-chan code.Identifier, str Storer) {
	for i := range identc {
		err := str.Save(i)
		if err != nil {
			// TODO
			log.Println(err)
		}
	}
}
