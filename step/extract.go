package step

import (
	"go/ast"
	"go/token"
	"reflect"

	"github.com/eroatta/src-reader/code"
)

const (
	typeFuncDecl = "FuncDecl"
	typeVarDecl  = "VarDecl"
)

// Extractor is used to define a custom identifier extractor.
type Extractor interface {
	// NodeType specifies the type of node that can be processed by the extractor.
	NodeType() reflect.Type
	// Extract retrieves the
	Extract(filename string, node ast.Node) (code.Identifier, bool)
}

// Extract TODO
// * Function declarations
// * Package variable/constant declarations
// * Method variable declarations
//		- range/loop declarations
// * Struct definitions
// * Interface definitions
// * Type definitions
func Extract(files []code.File) chan code.Identifier {
	extractors := make(map[reflect.Type]Extractor)

	fdcl := funcDeclExtractor{nil}
	extractors[fdcl.NodeType()] = fdcl

	identc := make(chan code.Identifier)
	go func() {
		for _, f := range files {
			if f.AST == nil {
				continue
			}

			ast.Inspect(f.AST, func(n ast.Node) bool {
				extractor := extractors[reflect.TypeOf(n)]
				if extractor == nil {
					return true
				}

				//log.Println("Type for: ", reflect.TypeOf(n))
				if ident, ok := extractor.Extract(f.Name, n); ok {
					identc <- ident
				}

				return true
			})
		}

		close(identc)
	}()

	return identc
}

func newIdent(name string, filename string, position token.Pos, declType string) code.Identifier {
	return code.Identifier{
		File:       filename,
		Position:   position,
		Name:       name,
		Type:       declType,
		Splits:     make(map[string][]string),
		Expansions: make(map[string][]string),
	}
}

type funcDeclExtractor struct {
	node *ast.FuncDecl
}

func (e funcDeclExtractor) NodeType() reflect.Type {
	return reflect.TypeOf(e.node)
}

func (e funcDeclExtractor) Extract(filename string, node ast.Node) (code.Identifier, bool) {
	fn, ok := node.(*ast.FuncDecl)
	if ok {
		return newIdent(fn.Name.String(), filename, fn.Pos(), typeFuncDecl), ok
	}

	return code.Identifier{}, false
}
