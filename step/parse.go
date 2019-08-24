package step

import (
	"go/parser"
	"go/token"
	"log"

	"github.com/eroatta/src-reader/code"
)

// Parse parses a file and creates an Abstract Syntax Tree (AST)
func Parse(filesc <-chan code.File) chan code.File {
	fset := token.NewFileSet()

	parsedc := make(chan code.File)
	go func() {
		for file := range filesc {
			node, err := parser.ParseFile(fset, file.Name, file.Raw, parser.ParseComments)
			if err != nil {
				// TODO: do we need error channel?
				log.Fatal(err)
				continue
			}

			file.AST = node
			file.FileSet = fset

			parsedc <- file
		}
		close(parsedc)
	}()

	return parsedc
}
