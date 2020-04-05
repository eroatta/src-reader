package step

import (
	"go/parser"
	"go/token"

	"github.com/eroatta/src-reader/entity"
)

// Parse parses a file and creates an Abstract Syntax Tree (AST) representation.
// It handles and returns a channel of entity.File elements.
func Parse(filesc <-chan entity.File) chan entity.File {
	fset := token.NewFileSet()

	parsedc := make(chan entity.File)
	go func() {
		for file := range filesc {
			node, err := parser.ParseFile(fset, file.Name, file.Raw, parser.ParseComments)

			file.AST = node
			file.FileSet = fset
			file.Error = err
			parsedc <- file
		}

		close(parsedc)
	}()

	return parsedc
}

// Merge joins files when necessary.
func Merge(parsedc <-chan entity.File) []entity.File {
	files := make([]entity.File, 0)
	for file := range parsedc {
		files = append(files, file)
	}

	return files
}
