package code

import (
	"go/ast"
	"go/token"
)

// File TODO
type File struct {
	Name    string
	Raw     []byte
	AST     *ast.File
	FileSet *token.FileSet
}

// Identifier TODO
type Identifier struct {
	File       string
	Position   token.Pos
	Name       string
	Type       string
	Node       *ast.Node
	Splits     map[string][]string
	Expansions map[string][]string
}
