package processors

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"

	"github.com/eroatta/token-splitex/splitters"
)

// SplitOn TODO
func SplitOn(fs *token.FileSet, node ast.Node, splitters ...splitters.Splitter) {
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			//log.Println(fmt.Sprintf("Found FuncDecl for :%s", fn.Name))
			//ast.Print(fs, fn)
			for _, splitter := range splitters {
				splits, err := splitter.Split(fn.Name.String())
				if err != nil {
					log.Fatalf("Error while splitting token \"%s\": %v", fn.Name.String(), err)
				}

				log.Println(fmt.Sprintf("FuncDecl \"%s\" Splitted into: %v", fn.Name.String(), splits))
			}
		}

		return true
	})
}
