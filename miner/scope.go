package miner

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/eroatta/src-reader/entity"
)

// Scope represents a scopes miner, which extracts information about
// the scope for each function/variable/struct/interface declaration.
type Scope struct {
	Filename        string
	PackageName     string
	PackageComments []string
	Comments        []*ast.CommentGroup
	Included        []ast.Decl
	Scopes          map[string]ScopedDecl
}

// NewScope initializes a new scopes miner.
func NewScope(filename string) Scope {
	return Scope{
		Filename:        filename,
		Scopes:          make(map[string]ScopedDecl),
		PackageComments: make([]string, 0),
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

// Type returns the specific entity.MinerType for the miner.
func (m Scope) Type() entity.MinerType {
	return entity.ScopedDeclarations
}

// Visit implements the ast.Visitor interface and handles the logic for the data extraction.
func (m Scope) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch elem := node.(type) {
	case *ast.File:
		m.PackageName = elem.Name.String()
		if elem.Doc != nil {
			for _, comment := range elem.Doc.List {
				m.PackageComments = append(m.PackageComments, cleanComment(comment.Text))
			}
		}

		m.Included = elem.Decls
		m.Comments = append(m.Comments, elem.Comments...)

	case *ast.FuncDecl:
		name := elem.Name.String()
		funcScopedDecl := newScopedDecl(m.PackageName, name, token.FUNC)

		// inbound and outbound parameters as variable declarations
		variableDecls := make([]string, 0)
		for _, in := range elem.Type.Params.List {
			var paramType string
			switch pType := in.Type.(type) {
			case *ast.Ident:
				paramType = pType.Name
			case *ast.ArrayType:
				if ident, ok := pType.Elt.(*ast.Ident); ok {
					paramType = fmt.Sprintf("[]%s", ident.Name)
				}
			case *ast.StructType:
				paramType = "struct"
			case *ast.FuncType:
				paramType = "func"
			case *ast.InterfaceType:
				paramType = "interface"
			case *ast.MapType:
				paramType = "map"
			case *ast.ChanType:
				paramType = "chan"
			default:
				paramType = "unknown"
			}

			for _, arg := range in.Names {
				if arg.Name != "" {
					variableDecls = append(variableDecls, fmt.Sprintf("%s %s", arg.String(), paramType))
				}
			}
		}

		if elem.Type.Results != nil {
			for _, out := range elem.Type.Results.List {
				var paramType string
				switch pType := out.Type.(type) {
				case *ast.Ident:
					paramType = pType.Name
				case *ast.ArrayType:
					if ident, ok := pType.Elt.(*ast.Ident); ok {
						paramType = fmt.Sprintf("[]%s", ident.Name)
					}
				case *ast.StructType:
					paramType = "struct"
				case *ast.FuncType:
					paramType = "func"
				case *ast.InterfaceType:
					paramType = "interface"
				case *ast.MapType:
					paramType = "map"
				case *ast.ChanType:
					paramType = "chan"
				default:
					paramType = "unknown"
				}

				for _, arg := range out.Names {
					if arg.Name != "" {
						variableDecls = append(variableDecls, fmt.Sprintf("%s %s", arg.String(), paramType))
					}
				}
			}
		}
		funcScopedDecl.VariableDecls = variableDecls

		// comments as doc for the function decl
		comments := make([]string, 0)
		if elem.Doc != nil {
			for _, comment := range elem.Doc.List {
				comments = append(comments, cleanComment(comment.Text))
			}
		}

		// comments inside the function decl
		start, end := elem.Pos(), elem.End()
		for _, group := range m.Comments {
			for _, comment := range group.List {
				if comment.Slash > start && comment.Slash < end {
					comments = append(comments, cleanComment(comment.Text))
				}
			}
		}

		funcScopedDecl.Comments = comments
		funcScopedDecl.PackageComments = m.PackageComments

		m.Scopes[funcScopedDecl.ID] = funcScopedDecl

	case *ast.GenDecl:
		if !(elem.Tok == token.VAR || elem.Tok == token.CONST || elem.Tok == token.TYPE) {
			return m
		}

		if !m.shouldMine(elem) {
			return m
		}

		comments := make([]string, 0)
		if elem.Doc != nil {
			for _, comment := range elem.Doc.List {
				comments = append(comments, cleanComment(comment.Text))
			}
		}

		for _, spec := range elem.Specs {
			if valSpec, ok := spec.(*ast.ValueSpec); ok {
				if valSpec.Doc != nil {
					for _, comment := range valSpec.Doc.List {
						comments = append(comments, cleanComment(comment.Text))
					}
				}

				if valSpec.Comment != nil {
					for _, comment := range valSpec.Comment.List {
						comments = append(comments, cleanComment(comment.Text))
					}
				}

				for j, name := range valSpec.Names {
					if name.Name == "_" {
						continue
					}

					varScopedDecl := newScopedDecl(m.PackageName, name.String(), elem.Tok)
					if valSpec.Values != nil {
						if val, ok := valSpec.Values[j].(*ast.BasicLit); ok && val.Kind == token.STRING {
							valStr := strings.Replace(val.Value, "\"", "", -1)
							varScopedDecl.BodyText = append(varScopedDecl.BodyText, cleanComment(valStr))
						}
					}
					varScopedDecl.Comments = comments
					varScopedDecl.PackageComments = m.PackageComments

					m.Scopes[varScopedDecl.ID] = varScopedDecl
				}
			}

			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				specificComments := make([]string, 0)
				if typeSpec.Doc != nil {
					for _, comment := range typeSpec.Doc.List {
						specificComments = append(specificComments, cleanComment(comment.Text))
					}
				}

				name := typeSpec.Name.String()

				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					structScopedDecl := newScopedDecl(m.PackageName, name, token.STRUCT)

					if structType.Fields != nil && structType.Fields.List != nil {
						variableDecls := make([]string, 0)
						for _, field := range structType.Fields.List {
							if fieldType, ok := field.Type.(*ast.Ident); ok {
								for _, name := range field.Names {
									variableDecls = append(variableDecls, fmt.Sprintf("%s %s", name.String(), fieldType.String()))
								}

								if field.Doc != nil {
									for _, comment := range field.Doc.List {
										specificComments = append(specificComments, cleanComment(comment.Text))
									}
								}
							}
						}
						structScopedDecl.VariableDecls = variableDecls
					}
					structScopedDecl.Comments = append(comments, specificComments...)
					structScopedDecl.PackageComments = m.PackageComments

					m.Scopes[structScopedDecl.ID] = structScopedDecl
				}

				if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
					interfaceScopedDecl := newScopedDecl(m.PackageName, name, token.INTERFACE)

					if interfaceType.Methods != nil && interfaceType.Methods.List != nil {
						statements := make([]string, 0)
						for _, method := range interfaceType.Methods.List {
							for _, name := range method.Names {
								statements = append(statements, name.String())
							}

							if method.Doc != nil {
								for _, comment := range method.Doc.List {
									specificComments = append(specificComments, cleanComment(comment.Text))
								}
							}
						}
						interfaceScopedDecl.Statements = statements
					}
					interfaceScopedDecl.Comments = append(comments, specificComments...)
					interfaceScopedDecl.PackageComments = m.PackageComments

					m.Scopes[interfaceScopedDecl.ID] = interfaceScopedDecl
				}
			}
		}
	}

	return m
}

func (m Scope) shouldMine(elem *ast.GenDecl) bool {
	var shouldMine bool
	for _, d := range m.Included {
		if d.Pos() == elem.Pos() {
			shouldMine = true
			break
		}
	}

	return shouldMine
}

// ScopedDeclarations returns a map of declaration IDs and the mined scope for each declaration.
func (m Scope) ScopedDeclarations() map[string]ScopedDecl {
	return m.Scopes
}
