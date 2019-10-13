package code

import (
	"go/ast"
	"go/token"
)

// Repository holds information of a GitHub repository.
type Repository struct {
	Name        string
	URL         string
	Hash        string
	DateCreated string
	Error       error
}

// File TODO
type File struct {
	Name    string
	Raw     []byte
	AST     *ast.File
	FileSet *token.FileSet
	Error   error
}

// Identifier TODO
type Identifier struct {
	File       string
	Position   token.Pos
	Name       string
	Type       string
	Parent     string
	ParentPos  token.Pos
	Node       *ast.Node
	Splits     map[string][]string
	Expansions map[string][]string
	Error      error
}
