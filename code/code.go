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
