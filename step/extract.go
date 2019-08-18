package step

import (
	"go/ast"

	"github.com/eroatta/src-reader/code"
)

// Extract TODO
func Extract(files []code.File) chan code.Identifier {
	identc := make(chan code.Identifier)
	go func() {
		for _, f := range files {
			ast.Inspect(f.AST, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if ok {
					ident := code.Identifier{
						File:       f.Name,
						Position:   fn.Pos(),
						Name:       fn.Name.String(),
						Type:       "FuncDecl",
						Splits:     make(map[string][]string),
						Expansions: make(map[string][]string),
					}
					identc <- ident
				}

				return true
			})
		}
		close(identc)
	}()

	return identc
}
