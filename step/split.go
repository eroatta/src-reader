package step

import (
	"log"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/token/conserv"
	"github.com/eroatta/token/greedy"
	"github.com/eroatta/token/lists"
	"github.com/eroatta/token/samurai"
)

// Split TODO
func Split(identc <-chan code.Identifier, tCtx samurai.TokenContext) chan code.Identifier {
	/*
		SPLITTING STEP: it splits all the identifiers in a AST, applying the given set of Splitters.
		Input: a list of *ast.File nodes, and a list of Splitters.
		Output: ?
	*/
	splittedc := make(chan code.Identifier)
	go func() {
		for ident := range identc {
			// TODO: remove logging
			log.Println("Splitting...")
			ident.Splits["conserv"] = conserv.Split(ident.Name)
			ident.Splits["greedy"] = greedy.Split(ident.Name, greedy.DefaultList)
			ident.Splits["samurai"] = samurai.Split(ident.Name, tCtx, lists.Prefixes, lists.Suffixes)

			splittedc <- ident
		}
		close(splittedc)
	}()

	return splittedc
}
