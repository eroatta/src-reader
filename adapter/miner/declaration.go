package miner

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/eroatta/nounphrases"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/conserv"
	"github.com/eroatta/token/lists"
	log "github.com/sirupsen/logrus"
)

// NewDeclarationsFactory creates a new declarations miner factory.
func NewDeclarationsFactory() entity.MinerFactory {
	return declarationsFactory{}
}

type declarationsFactory struct{}

func (f declarationsFactory) Make() (entity.Miner, error) {
	return NewDeclaration(lists.Dictionary), nil
}

// NewDeclaration initializes a new declarations miner.
func NewDeclaration(dict lists.List) *Declaration {
	return &Declaration{
		miner: miner{"declarations"},
		Dict:  dict,
		Decls: make(map[string]entity.Decl),
	}
}

// Declaration represents the declarations miner, which extracts information about
// words and phrases for each function/variable/struct/interface declaration.
type Declaration struct {
	miner
	Filename    string
	Dict        lists.List
	PackageName string
	Comments    []*ast.CommentGroup
	Included    []ast.Decl
	Decls       map[string]entity.Decl
}

// SetCurrentFile specifies the current file being mined.
func (m *Declaration) SetCurrentFile(filename string) {
	m.Filename = filename
}

// Visit implements the ast.Visitor interface and handles the logic for the data extraction.
func (m *Declaration) Visit(node ast.Node) ast.Visitor {
	// TODO: (stemming + stopping)
	if node == nil {
		return nil
	}

	switch elem := node.(type) {
	case *ast.File:
		m.PackageName = elem.Name.String()
		m.Included = elem.Decls
		m.Comments = append(m.Comments, elem.Comments...)

	case *ast.FuncDecl:
		functionText := extractDeclFromFunction(elem, m)
		m.Decls[functionText.ID] = functionText

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
				genDeclComments = extractWordAndPhrasesFromComment(genDeclComments, comment.Text, m.Dict)
			}
		}

		for _, spec := range elem.Specs {
			if valSpec, ok := spec.(*ast.ValueSpec); ok {
				valDeclComments := newDecl("val", elem.Tok)
				if valSpec.Doc != nil {
					for _, comment := range valSpec.Doc.List {
						valDeclComments = extractWordAndPhrasesFromComment(valDeclComments, comment.Text, m.Dict)
					}
				}

				for j, name := range valSpec.Names {
					if name.Name == "_" {
						continue
					}

					declText := newDecl(declID(m.Filename, m.PackageName, elem.Tok, name.String(), ""), elem.Tok)
					declText = extractDeclFromValue(declText, valSpec, name.Name, j, m.Dict)
					declText = merge(merge(declText, genDeclComments), valDeclComments)

					m.Decls[declText.ID] = declText
				}
			}

			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				name := typeSpec.Name.String()
				declText := newDecl("", token.TYPE)

				for _, part := range strings.Split(conserv.Split(name), " ") {
					if m.Dict.Contains(part) {
						declText.Words[part] = struct{}{}
					}
				}

				if typeSpec.Doc != nil {
					for _, comment := range typeSpec.Doc.List {
						declText = extractWordAndPhrasesFromComment(declText, comment.Text, m.Dict)
					}
				}

				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					declText.ID = declID(m.Filename, m.PackageName, token.STRUCT, name, "")
					declText.DeclType = token.STRUCT
					declText = extractDeclFromStruct(declText, structType, m.Dict)
				}

				if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
					declText.ID = declID(m.Filename, m.PackageName, token.INTERFACE, name, "")
					declText.DeclType = token.INTERFACE
					declText = extractDeclFromInterface(declText, interfaceType, m.Dict)
				}

				declText = merge(declText, genDeclComments)

				m.Decls[declText.ID] = declText
			}
		}
	}

	return m
}

func newDecl(ID string, declType token.Token) entity.Decl {
	return entity.Decl{
		ID:       ID,
		DeclType: declType,
		Words:    make(map[string]struct{}),
		Phrases:  make(map[string]struct{}),
	}
}

func (m Declaration) shouldMine(elem *ast.GenDecl) bool {
	var shouldMine bool
	for _, d := range m.Included {
		if d.Pos() == elem.Pos() {
			shouldMine = true
			break
		}
	}

	return shouldMine
}

func extractDeclFromFunction(elem *ast.FuncDecl, m *Declaration) entity.Decl {
	name := elem.Name.String()
	receiver := ""
	if elem.Recv != nil && elem.Recv.NumFields() > 0 {
		for _, r := range elem.Recv.List {
			typ, ok := r.Type.(*ast.Ident)
			if ok {
				receiver = typ.Name
			}
		}
	}

	functionText := newDecl(declID(m.Filename, m.PackageName, token.FUNC, name, receiver), token.FUNC)

	for _, part := range strings.Split(conserv.Split(name), " ") {
		if m.Dict.Contains(part) {
			functionText.Words[part] = struct{}{}
		}
	}

	if elem.Doc != nil {
		for _, comment := range elem.Doc.List {
			functionText = extractWordAndPhrasesFromComment(functionText, comment.Text, m.Dict)
		}
	}

	start, end := elem.Pos(), elem.End()
	for _, group := range m.Comments {
		for _, comment := range group.List {
			if comment.Slash > start && comment.Slash < end {
				functionText = extractWordAndPhrasesFromComment(functionText, comment.Text, m.Dict)
			}
		}
	}

	return functionText
}

func extractDeclFromValue(declText entity.Decl, valSpec *ast.ValueSpec, name string, index int, list lists.List) entity.Decl {
	for _, part := range strings.Split(conserv.Split(name), " ") {
		if list.Contains(part) {
			declText.Words[part] = struct{}{}
		}
	}

	if valSpec.Values != nil {
		if val, ok := valSpec.Values[index].(*ast.BasicLit); ok && val.Kind == token.STRING {
			valStr := strings.Replace(val.Value, "\"", "", -1)
			for _, word := range strings.Split(valStr, " ") {
				word = strings.ToLower(cleaner.ReplaceAllString(word, ""))
				if list.Contains(word) {
					declText.Words[word] = struct{}{}
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

func extractDeclFromStruct(declText entity.Decl, structType *ast.StructType, list lists.List) entity.Decl {
	if structType.Fields != nil && structType.Fields.List != nil {
		for _, field := range structType.Fields.List {
			for _, fname := range field.Names {
				for _, part := range strings.Split(conserv.Split(fname.Name), " ") {
					if list.Contains(part) {
						declText.Words[part] = struct{}{}
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

func extractDeclFromInterface(declText entity.Decl, interfaceType *ast.InterfaceType, list lists.List) entity.Decl {
	if interfaceType.Methods != nil && interfaceType.Methods.List != nil {
		for _, method := range interfaceType.Methods.List {
			for _, mname := range method.Names {
				for _, part := range strings.Split(conserv.Split(mname.Name), " ") {
					if list.Contains(part) {
						declText.Words[part] = struct{}{}
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

func merge(a entity.Decl, b entity.Decl) entity.Decl {
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
	cleanComment = strings.ToLower(cleanComment)

	return strings.TrimSpace(cleanComment)
}

func extractWordAndPhrasesFromComment(functionText entity.Decl, comment string, list lists.List) entity.Decl {
	cleanComment := cleanComment(comment)
	for _, word := range strings.Split(cleanComment, " ") {
		word = cleaner.ReplaceAllString(word, "")
		if list.Contains(word) {
			functionText.Words[word] = struct{}{}
		}
	}

	phrases, err := nounphrases.Find(cleanComment)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("unable to retrieve phrases from comment \"%s\"", cleanComment))
	}
	for _, phr := range phrases {
		functionText.Phrases[phr] = struct{}{}
	}

	return functionText
}

func declID(filename string, pkg string, declType token.Token, name string, receiver string) string {
	return entity.NewDeclarationIDBuilder().
		WithFilename(filename).
		WithPackage(pkg).
		WithReceiver(receiver).
		WithName(name).
		WithType(declType).
		Build()
}

// Results returns a map of declaration IDs and the mined text for each declaration.
func (m Declaration) Results() interface{} {
	return m.Decls
}
