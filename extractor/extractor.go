package extractor

import (
	"go/ast"
	"go/token"

	"github.com/eroatta/src-reader/code"
)

var types = map[token.Token]string{
	token.CONST:     "ConstDecl",
	token.INTERFACE: "InterfaceDecl",
	token.STRUCT:    "StructDecl",
	token.VAR:       "VarDecl",
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
				e.identifiers = append(e.identifiers, fromValueSpec(e.filename, elem.Tok, decl)...)
			case *ast.TypeSpec:
				e.identifiers = append(e.identifiers, fromTypeSpec(e.filename, decl)...)
			}
		}
	}

	return e
}

func fromValueSpec(filename string, token token.Token, decl *ast.ValueSpec) []code.Identifier {
	identifiers := []code.Identifier{}
	for _, name := range decl.Names {
		if name.Name == "_" {
			continue
		}

		identifier := code.Identifier{
			File:       filename,
			Position:   name.Pos(),
			Name:       name.String(),
			Type:       types[token],
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		}

		identifiers = append(identifiers, identifier)
	}

	return identifiers
}

func fromTypeSpec(filename string, decl *ast.TypeSpec) []code.Identifier {
	var identifierType string
	switch decl.Type.(type) {
	case *ast.StructType:
		identifierType = types[token.STRUCT]
	case *ast.InterfaceType:
		identifierType = types[token.INTERFACE]
	default:
		return []code.Identifier{}
	}

	identifiers := []code.Identifier{
		newIdentifier(filename, decl.Name.Pos(), decl.Name.String(), identifierType),
	}

	return identifiers
}

func newIdentifier(filename string, pos token.Pos, name string, identifierType string) code.Identifier {
	return code.Identifier{
		File:       filename,
		Position:   pos,
		Name:       name,
		Type:       identifierType,
		Splits:     make(map[string][]string),
		Expansions: make(map[string][]string),
	}
}

// Identifiers returns the list of found identifiers.
func (e *Extractor) Identifiers() []code.Identifier {
	return e.identifiers
}
