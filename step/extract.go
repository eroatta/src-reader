package step

import (
	"go/ast"
	"reflect"

	"github.com/eroatta/src-reader/code"
)

// Extractor is used to define a custom identifier extractor.
type Extractor interface {
	// NodeType specifies the type of node that can be processed by the extractor.
	NodeType() reflect.Type
	// Extract retrieves an array of identifiers found on the node processed by the extractor.
	Extract(filename string, node ast.Node) []code.Identifier
}

// Extract traverses each Abstract Syntax Tree and applies a set of extractors
// to retrieve the identifiers that are interest of us.
func Extract(files []code.File, extractors ...Extractor) chan code.Identifier {
	mappedExtractors := make(map[reflect.Type]Extractor)
	for _, ext := range extractors {
		mappedExtractors[ext.NodeType()] = ext
	}

	identc := make(chan code.Identifier)
	go func() {
		for _, f := range files {
			if f.AST == nil {
				continue
			}

			ast.Inspect(f.AST, func(n ast.Node) bool {
				if extractor, ok := mappedExtractors[reflect.TypeOf(n)]; ok {
					for _, ident := range extractor.Extract(f.Name, n) {
						identc <- ident
					}
				}

				return true
			})
		}

		close(identc)
	}()

	return identc
}
