package extractor

import (
	"go/ast"
	"go/token"

	"github.com/eroatta/src-reader/entity"
)

// Extractor represents an extraction algorithm, capable of retriving several definitions/declarations.
// The current implementation handles:
//	* function declaration
//	* variable declaration
//	* constant declaration
//	* struct definition
//	* interface definition
type Extractor struct {
	filename      string
	packageName   string
	currentLoc    string
	currentLocPos token.Pos
	identifiers   []entity.Identifier
	scopes        map[string]*ast.Object
}

// New creates a new Extractor.
func New(filename string) entity.Extractor {
	return &Extractor{
		filename:    filename,
		identifiers: make([]entity.Identifier, 0),
	}
}

// Visit implements the ast.Visitor interface and handles the logic for the identifiers extraction.
func (e *Extractor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch elem := node.(type) {
	case *ast.File:
		e.packageName = elem.Name.String()
		e.scopes = elem.Scope.Objects

	case *ast.FuncDecl:
		name := elem.Name.String()

		recv := ""
		if elem.Recv != nil && elem.Recv.NumFields() > 0 {
			for _, r := range elem.Recv.List {
				switch exp := r.Type.(type) {
				case *ast.Ident:
					recv = exp.Name
				case *ast.StarExpr:
					typ, ok := exp.X.(*ast.Ident)
					if ok {
						recv = typ.Name
					}
				}
			}
		}

		id := entity.NewIDBuilder().WithFilename(e.filename).
			WithPackage(e.packageName).WithReceiver(recv).WithName(name).WithType(token.FUNC).Build()
		e.identifiers = append(e.identifiers, newIdentifier(id, e.filename, elem.Pos(), name, token.FUNC))

		// set current location at the beginning of each function
		e.currentLoc = id
		e.currentLocPos = elem.Pos()

	case *ast.GenDecl:
		for _, spec := range elem.Specs {
			switch decl := spec.(type) {
			case *ast.ValueSpec:
				e.identifiers = append(e.identifiers, e.fromValueSpec(e.filename, elem.Tok, decl)...)
			case *ast.TypeSpec:
				e.identifiers = append(e.identifiers, e.fromTypeSpec(e.filename, decl)...)
			}
		}
	}

	return e
}

func (e *Extractor) fromValueSpec(filename string, token token.Token, decl *ast.ValueSpec) []entity.Identifier {
	identifiers := []entity.Identifier{}
	for _, name := range decl.Names {
		if name.Name == "_" {
			continue
		}

		id := entity.NewIDBuilder().WithFilename(e.filename).
			WithPackage(e.packageName).WithName(name.String()).WithType(token).Build()

		identifiers = append(identifiers,
			newIdentifier(id, filename, name.Pos(), name.String(), token))
	}

	return identifiers
}

func (e *Extractor) fromTypeSpec(filename string, decl *ast.TypeSpec) []entity.Identifier {
	var identifierType token.Token
	switch decl.Type.(type) {
	case *ast.StructType:
		identifierType = token.STRUCT
	case *ast.InterfaceType:
		identifierType = token.INTERFACE
	default:
		return []entity.Identifier{}
	}

	id := entity.NewIDBuilder().WithFilename(e.filename).
		WithPackage(e.packageName).WithName(decl.Name.String()).WithType(identifierType).Build()

	return []entity.Identifier{
		newIdentifier(id, filename, decl.Pos(), decl.Name.String(), identifierType)}
}

func newIdentifier(id string, filename string, pos token.Pos, name string, identifierType token.Token) entity.Identifier {
	return entity.Identifier{
		ID:         id,
		File:       filename,
		Position:   pos,
		Name:       name,
		Type:       identifierType,
		Splits:     make(map[string][]entity.Split),
		Expansions: make(map[string][]entity.Expansion),
	}
}

// Identifiers returns the list of found identifiers.
func (e *Extractor) Identifiers() []entity.Identifier {
	return e.identifiers
}
