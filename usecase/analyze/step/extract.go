package step

import (
	"go/ast"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
)

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
