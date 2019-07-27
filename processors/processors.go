package processors

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
)

// SplitOn TODO
func SplitOn(fs *token.FileSet, node ast.Node, splitters ...func(string) []string) {
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			//log.Println(fmt.Sprintf("Found FuncDecl for :%s", fn.Name))
			//ast.Print(fs, fn)
			for _, splitFunc := range splitters {
				splits := splitFunc(fn.Name.String())

				log.Println(fmt.Sprintf("FuncDecl \"%s\" Splitted into: %v", fn.Name.String(), splits))
			}
		}

		return true
	})
}
