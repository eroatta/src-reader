package miner

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/eroatta/token/amap"
)

// Scope represents a scopes miner, which extracts information about
// the scope for each function/variable/struct/interface declaration.
type Scope struct {
	filename        string
	packageName     string
	packageComments []string
	comments        []*ast.CommentGroup
	included        []ast.Decl
	scopes          map[string]ScopedDecl
}

// NewScope initializes a new scopes miner.
func NewScope(filename string) Scope {
	return Scope{
		filename:        filename,
		scopes:          make(map[string]ScopedDecl),
		packageComments: make([]string, 0),
	}
}

// ScopedDecl represents the related scope for a declaration.
type ScopedDecl struct {
	ID              string
	DeclType        token.Token
	Name            string
	VariableDecls   []string
	Statements      []string
	BodyText        []string
	Comments        []string
	PackageComments []string
}

func newScopedDecl(pkg string, name string, declType token.Token) ScopedDecl {
	return ScopedDecl{
		ID:              declID(pkg, declType, name),
		DeclType:        declType,
		Name:            name,
		VariableDecls:   make([]string, 0),
		Statements:      make([]string, 0),
		BodyText:        make([]string, 0),
		Comments:        make([]string, 0),
		PackageComments: make([]string, 0),
	}
}

// Name returns the specific name for the miner.
func (m Scope) Name() string {
	return "scope"
}

// Visit implements the ast.Visitor interface and handles the logic for the data extraction.
func (m Scope) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch elem := node.(type) {
	case *ast.File:
		m.packageName = elem.Name.String()
		if elem.Doc != nil {
			for _, comment := range elem.Doc.List {
				m.packageComments = append(m.packageComments, cleanComment(comment.Text))
			}
		}

		m.included = elem.Decls
		m.comments = append(m.comments, elem.Comments...)

	case *ast.FuncDecl:
		name := elem.Name.String()
		funcScopedDecl := newScopedDecl(m.packageName, name, token.FUNC)

		// set function comments
		if elem.Doc != nil {
			for _, comment := range elem.Doc.List {
				funcScopedDecl.Comments = append(funcScopedDecl.Comments, cleanComment(comment.Text))
			}
		}

		// set package comments
		funcScopedDecl.PackageComments = m.packageComments
		m.scopes[funcScopedDecl.ID] = funcScopedDecl

	case *ast.GenDecl:
		fmt.Println(elem)
	}

	return m
}

// Amap represents an AMAP miner, which extracts the scopes from a file.
type Amap struct {
	name         string
	file         string
	fileComments []string
	scopes       map[string]amap.TokenScope
}

// Visit implements the ast.Visitor interface and handles the logic the node mining.
func (m Amap) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	// TODO: what about package variables?

	switch elem := node.(type) {
	case *ast.File:
		for _, commentGroup := range elem.Comments {
			for _, comment := range commentGroup.List {
				m.fileComments = append(m.fileComments, cleanComment(comment.Text))
			}
		}

	case *ast.FuncDecl:
		name := elem.Name.String()

		funcComments := make([]string, 0)
		if elem.Doc != nil {
			for _, comment := range elem.Doc.List {
				funcComments = append(funcComments, cleanComment(comment.Text))
			}
		}

		decls := make([]string, 0)
		for _, in := range elem.Type.Params.List {
			typ, ok := in.Type.(*ast.Ident)
			if !ok {
				continue
			}
			for _, arg := range in.Names {
				decls = append(decls, fmt.Sprintf("%s %s", typ.Name, arg.String()))
			}
		}

		if elem.Type.Results != nil {
			for _, out := range elem.Type.Results.List {
				typ, ok := out.Type.(*ast.Ident)
				if !ok {
					continue
				}
				for _, arg := range out.Names {
					if arg.Name != "" {
						decls = append(decls, fmt.Sprintf("%s %s", typ.Name, arg.Name))
					}
				}
			}
		}

		// replace by []string on token project
		bodyText := ""
		// TODO: add body text extraction

		tokenScope := amap.NewTokenScope(decls, name, bodyText, funcComments, m.fileComments)
		key := fmt.Sprintf("%s+%s", m.file, name)
		m.scopes[key] = tokenScope
	}

	return m
}

// ScopedDeclarations returns a map of declaration IDs and the mined scope for each declaration.
func (m Scope) ScopedDeclarations() map[string]ScopedDecl {
	return m.scopes
}
