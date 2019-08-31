package extractor

import (
	"go/ast"
	"reflect"

	"github.com/eroatta/src-reader/code"
)

// * Function declarations
// * Package variable/constant declarations
// * Method variable declarations
//		- range/loop declarations
// * Struct definitions
// * Interface definitions
// * Type definitions

// FuncDecl performs the extraction of identifiers on a function declaration node.
type FuncDecl struct {
	node *ast.FuncDecl
}

// NodeType returns the node handled by FuncDeclExtractor.
func (e FuncDecl) NodeType() reflect.Type {
	return reflect.TypeOf(e.node)
}

// Extract TODO
func (e FuncDecl) Extract(filename string, node ast.Node) []code.Identifier {
	fn, ok := node.(*ast.FuncDecl)
	if ok {
		return []code.Identifier{
			{
				File:       filename,
				Position:   fn.Name.Pos(),
				Name:       fn.Name.String(),
				Type:       "FuncDecl",
				Splits:     make(map[string][]string),
				Expansions: make(map[string][]string),
			}}
	}

	return []code.Identifier{}
}
