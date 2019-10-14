package step

import (
	"go/ast"

	"github.com/eroatta/src-reader/code"
)

// ExtractorFactoryFunc defines a function to create a new Extractor.
type ExtractorFactoryFunc func(filename string) Extractor

// Extractor is used to define a custom identifier extractor.
type Extractor interface {
	// Visit applies the extraction logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
	// Identifiers returns the extracted identifiers.
	Identifiers() []code.Identifier
}

// Extract traverses each Abstract Syntax Tree and applies a set of extractors
// to retrieve the identifiers that are interest of us.
func Extract(files []code.File, ef ExtractorFactoryFunc) chan code.Identifier {
	identc := make(chan code.Identifier)
	go func() {
		for _, f := range files {
			if f.AST == nil {
				continue
			}

			e := ef(f.Name)
			ast.Walk(e, f.AST)

			for _, ident := range e.Identifiers() {
				identc <- ident
			}
		}

		close(identc)
	}()

	return identc
}
