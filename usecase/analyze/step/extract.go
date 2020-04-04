package step

import (
	"go/ast"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
)

// TODO: remove
// ExtractorFactoryFunc defines a function to create a new Extractor.
//type ExtractorFactoryFunc func(filename string) Extractor

// Extractor is used to define a custom identifier extractor.
/*type Extractor interface {
	// Visit applies the extraction logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
	// Identifiers returns the extracted identifiers.
	Identifiers() []code.Identifier
}*/

// extract traverses each Abstract Syntax Tree and applies an extractor
// to retrieve the identifiers that are interest of us.
func Extract(files []code.File, factory entity.ExtractorFactory) chan code.Identifier {
	identc := make(chan code.Identifier)
	go func() {
		for _, f := range files {
			if f.AST == nil {
				continue
			}

			extractor := factory(f.Name)
			ast.Walk(extractor, f.AST)

			for _, ident := range extractor.Identifiers() {
				identc <- ident
			}
		}

		close(identc)
	}()

	return identc
}
