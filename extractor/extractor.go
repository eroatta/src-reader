package extractor

import (
	"go/ast"
	"go/token"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/step"
)

var types = map[token.Token]string{
	token.CONST:     "ConstDecl",
	token.DEFINE:    "AssignStmt",
	token.FUNC:      "FuncDecl",
	token.INTERFACE: "InterfaceDecl",
	token.STRUCT:    "StructDecl",
	token.VAR:       "VarDecl",
}

type Extractor struct {
	filename      string
	currentLoc    string
	currentLocPos token.Pos
	identifiers   []code.Identifier
}

// New creates a new Extractor.
func New(filename string) step.Extractor {
	return &Extractor{
		filename:    filename,
		identifiers: make([]code.Identifier, 0),
	}
}

// TODO: change
func New2(filename string) entity.Extractor {
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
	case *ast.FuncDecl:
		name := elem.Name.String()
		e.currentLoc = name
		e.currentLocPos = elem.Pos()

		e.identifiers = append(e.identifiers, newIdentifier(e.filename, elem.Pos(), name, types[token.FUNC]))

	case *ast.AssignStmt:
		if elem.Tok != token.DEFINE {
			return e
		}

		for _, name := range elem.Lhs {
			ident, ok := name.(*ast.Ident)
			if !ok {
				continue
			}

			if ident.Name == "_" || ident.Name == "" {
				continue
			}

			if ident.Obj != nil && ident.Obj.Pos() == ident.Pos() {
				e.identifiers = append(e.identifiers,
					newChildIdentifier(e.filename, ident.Pos(), ident.Name, types[elem.Tok], e.currentLoc, e.currentLocPos))
			}
		}

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

		identifiers = append(identifiers,
			newIdentifier(filename, name.Pos(), name.String(), types[token]))
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

func newChildIdentifier(filename string, pos token.Pos, name string, identifierType string, parent string, parentPos token.Pos) code.Identifier {
	i := newIdentifier(filename, pos, name, identifierType)
	i.Parent = parent
	i.ParentPos = parentPos

	return i
}

// Identifiers returns the list of found identifiers.
func (e *Extractor) Identifiers() []code.Identifier {
	return e.identifiers
}
