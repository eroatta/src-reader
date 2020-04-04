package storer

import (
	"fmt"
	"log"

	"github.com/eroatta/src-reader/code"
)

func New() interface{} {
	return stdout{}
}

type stdout struct {
}

func (m stdout) Save(ident code.Identifier) error {
	log.Println("Storing identifier...")
	for alg, splits := range ident.Splits {
		log.Println(fmt.Sprintf("%s \"%s\" Splitted into: %v by %s", ident.Type, ident.Name, splits, alg))
	}

	for alg, expans := range ident.Expansions {
		log.Println(fmt.Sprintf("%s \"%s\" Expanded into: %v by %s", ident.Type, ident.Name, expans, alg))
	}

	return nil
}
