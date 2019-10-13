package extractor

import (
	"go/ast"
	"go/token"

	"github.com/eroatta/src-reader/code"
)

var types = map[token.Token]string{
	token.CONST: "ConstDecl",
	token.VAR:   "VarDecl",
}

type Extractor struct {
	filename    string
	currentLoc  string
	identifiers []code.Identifier
}

// New creates a new Extractor.
func New(filename string) *Extractor {
	return &Extractor{
		filename:    filename,
		identifiers: make([]code.Identifier, 0),
	}
}

// Visit implements the ast.Visitor interface and handles the logic for the identifiers extraction.
func (e *Extractor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch elem := node.(type) {
	case *ast.GenDecl:
		for _, spec := range elem.Specs {
			switch decl := spec.(type) {
			case *ast.ValueSpec:
				for _, name := range decl.Names {
					if name.Name == "_" {
						continue
					}

					identifier := code.Identifier{
						// TODO, perhaps we can use just the position
						File:       e.filename,
						Position:   name.Pos(),
						Name:       name.String(),
						Type:       types[elem.Tok],
						Splits:     make(map[string][]string),
						Expansions: make(map[string][]string),
					}

					e.identifiers = append(e.identifiers, identifier)
				}
			}
		}
	}

	return e
}

// Identifiers returns the list of found identifiers.
func (e *Extractor) Identifiers() []code.Identifier {
	return e.identifiers
}
