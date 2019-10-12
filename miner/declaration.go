package miner

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/eroatta/nounphrases"
	"github.com/eroatta/token/conserv"
	"github.com/eroatta/token/lists"
)

// Decl contains the mined text (words and phrases) related to a declaration.
type Decl struct {
	ID       string
	DeclType token.Token
	Words    map[string]struct{}
	Phrases  map[string]struct{}
}

func newDecl(ID string, declType token.Token) Decl {
	return Decl{
		ID:       ID,
		DeclType: declType,
		Words:    make(map[string]struct{}),
		Phrases:  make(map[string]struct{}),
	}
}

// Declaration represents the declarations miner, which extracts information about
// words and phrases for each function/variable/struct/interface declaration.
type Declaration struct {
	dict        lists.List
	packageName string
	comments    []*ast.CommentGroup
	included    []ast.Decl
	decls       map[string]Decl
}

// NewDeclaration initializes a new declarations miner.
func NewDeclaration(dict lists.List) Declaration {
	return Declaration{
		dict:  dict,
		decls: make(map[string]Decl),
	}
}

// Name returns the specific name for the miner.
func (m Declaration) Name() string {
	return "declaration"
}

// Visit implements the ast.Visitor interface and handles the logic for the data extraction.
func (m Declaration) Visit(node ast.Node) ast.Visitor {
	// TODO: (stemming + stopping)
	if node == nil {
		return nil
	}

	switch elem := node.(type) {
	case *ast.File:
		m.packageName = elem.Name.String()
		m.included = elem.Decls
		m.comments = append(m.comments, elem.Comments...)

	case *ast.FuncDecl:
		functionText := extractDeclFromFunction(elem, m)
		m.decls[functionText.ID] = functionText

	case *ast.GenDecl:
		if elem.Tok != token.VAR && elem.Tok != token.CONST && elem.Tok != token.TYPE {
			return m
		}

		if !m.shouldMine(elem) {
			return m
		}

		genDeclComments := newDecl("common", elem.Tok)
		if elem.Doc != nil {
			for _, comment := range elem.Doc.List {
				genDeclComments = extractWordAndPhrasesFromComment(genDeclComments, comment.Text, m.dict)
			}
		}

		for _, spec := range elem.Specs {
			if valSpec, ok := spec.(*ast.ValueSpec); ok {
				valDeclComments := newDecl("val", elem.Tok)
				if valSpec.Doc != nil {
					for _, comment := range valSpec.Doc.List {
						valDeclComments = extractWordAndPhrasesFromComment(valDeclComments, comment.Text, m.dict)
					}
				}

				for j, name := range valSpec.Names {
					if name.Name == "_" {
						continue
					}

					declText := newDecl(declID(m.packageName, elem.Tok, name.String()), elem.Tok)
					declText = extractDeclFromValue(declText, valSpec, name.Name, j, m.dict)
					declText = merge(merge(declText, genDeclComments), valDeclComments)

					m.decls[declText.ID] = declText
				}
			}

			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				name := typeSpec.Name.String()
				declText := newDecl("", token.TYPE)

				for _, part := range conserv.Split(name) {
					if m.dict.Contains(part) {
						declText.Words[strings.ToLower(part)] = struct{}{}
					}
				}

				if typeSpec.Doc != nil {
					for _, comment := range typeSpec.Doc.List {
						declText = extractWordAndPhrasesFromComment(declText, comment.Text, m.dict)
					}
				}

				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					declText.ID = declID(m.packageName, token.STRUCT, name)
					declText.DeclType = token.STRUCT
					declText = extractDeclFromStruct(declText, structType, m.dict)
				}

				if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
					declText.ID = declID(m.packageName, token.INTERFACE, name)
					declText.DeclType = token.INTERFACE
					declText = extractDeclFromInterface(declText, interfaceType, m.dict)
				}

				declText = merge(declText, genDeclComments)

				m.decls[declText.ID] = declText
			}
		}
	}

	return m
}

func (m Declaration) shouldMine(elem *ast.GenDecl) bool {
	var shouldMine bool
	for _, d := range m.included {
		if d.Pos() == elem.Pos() {
			shouldMine = true
			break
		}
	}

	return shouldMine
}

func extractDeclFromFunction(elem *ast.FuncDecl, m Declaration) Decl {
	name := elem.Name.String()
	functionText := newDecl(declID(m.packageName, token.FUNC, name), token.FUNC)

	for _, part := range conserv.Split(name) {
		if m.dict.Contains(part) {
			functionText.Words[strings.ToLower(part)] = struct{}{}
		}
	}

	if elem.Doc != nil {
		for _, comment := range elem.Doc.List {
			functionText = extractWordAndPhrasesFromComment(functionText, comment.Text, m.dict)
		}
	}

	start, end := elem.Pos(), elem.End()
	for _, group := range m.comments {
		for _, comment := range group.List {
			if comment.Slash > start && comment.Slash < end {
				functionText = extractWordAndPhrasesFromComment(functionText, comment.Text, m.dict)
			}
		}
	}

	return functionText
}

func extractDeclFromValue(declText Decl, valSpec *ast.ValueSpec, name string, index int, list lists.List) Decl {
	for _, part := range conserv.Split(name) {
		if list.Contains(part) {
			declText.Words[strings.ToLower(part)] = struct{}{}
		}
	}

	if valSpec.Values != nil {
		if val, ok := valSpec.Values[index].(*ast.BasicLit); ok && val.Kind == token.STRING {
			valStr := strings.Replace(val.Value, "\"", "", -1)
			for _, word := range strings.Split(valStr, " ") {
				word = cleaner.ReplaceAllString(word, "")
				if list.Contains(strings.ToLower(word)) {
					declText.Words[strings.ToLower(word)] = struct{}{}
				}
			}

			phrases, _ := nounphrases.Find(cleanComment(valStr))
			for _, phr := range phrases {
				declText.Phrases[phr] = struct{}{}
			}
		}
	}

	return declText
}

func extractDeclFromStruct(declText Decl, structType *ast.StructType, list lists.List) Decl {
	if structType.Fields != nil && structType.Fields.List != nil {
		for _, field := range structType.Fields.List {
			for _, fname := range field.Names {
				for _, part := range conserv.Split(fname.Name) {
					if list.Contains(part) {
						declText.Words[strings.ToLower(part)] = struct{}{}
					}
				}
			}

			if field.Doc != nil {
				for _, comment := range field.Doc.List {
					declText = extractWordAndPhrasesFromComment(declText, comment.Text, list)
				}
			}
		}
	}

	return declText
}

func extractDeclFromInterface(declText Decl, interfaceType *ast.InterfaceType, list lists.List) Decl {
	if interfaceType.Methods != nil && interfaceType.Methods.List != nil {
		for _, method := range interfaceType.Methods.List {
			for _, mname := range method.Names {
				for _, part := range conserv.Split(mname.Name) {
					if list.Contains(part) {
						declText.Words[strings.ToLower(part)] = struct{}{}
					}
				}
			}

			if method.Doc != nil {
				for _, comment := range method.Doc.List {
					declText = extractWordAndPhrasesFromComment(declText, comment.Text, list)
				}
			}
		}
	}

	return declText
}

func merge(a Decl, b Decl) Decl {
	for k, v := range b.Words {
		a.Words[k] = v
	}

	for k, v := range b.Phrases {
		a.Phrases[k] = v
	}

	return a
}

func cleanComment(text string) string {
	cleanComment := strings.ReplaceAll(text, "//", "")
	cleanComment = strings.ReplaceAll(cleanComment, "\\n", " ")
	cleanComment = strings.ReplaceAll(cleanComment, "\\t", " ")

	return strings.TrimSpace(cleanComment)
}

func extractWordAndPhrasesFromComment(functionText Decl, comment string, list lists.List) Decl {
	cleanComment := cleanComment(comment)
	for _, word := range strings.Split(cleanComment, " ") {
		word = cleaner.ReplaceAllString(word, "")
		if list.Contains(word) {
			functionText.Words[strings.ToLower(word)] = struct{}{}
		}
	}

	if phrases, err := nounphrases.Find(cleanComment); err == nil {
		for _, phr := range phrases {
			functionText.Phrases[phr] = struct{}{}
		}
	} else {
		// TODO: improve logging
		fmt.Println(err)
	}

	return functionText
}

func declID(pkg string, declType token.Token, name string) string {
	return fmt.Sprintf("%s++%s::%s", pkg, declType, name)
}

// Decls returns a map of declaration IDs and the mined text for each declaration.
func (m Declaration) Decls() map[string]Decl {
	return m.decls
}