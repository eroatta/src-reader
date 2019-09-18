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
	// on a FuncDecl, get comments + identifiers + body text
	// and extract words (stemming + stopping)
	// and also extract phrases
	if node == nil {
		return nil
	}

	switch elem := node.(type) {
	case *ast.File:
		m.packageName = elem.Name.String()
		m.included = elem.Decls

		// collect inner comment to use on their proper function
		m.comments = append(m.comments, elem.Comments...)

	case *ast.FuncDecl:
		functionText := getFunctionTextFromFuncDecl(elem, m)
		m.decls[functionText.ID] = functionText

	case *ast.GenDecl:
		// TODO improve this code
		if elem.Tok != token.VAR && elem.Tok != token.CONST && elem.Tok != token.TYPE {
			return m
		}

		var weCare bool
		for _, d := range m.included {
			if d.Pos() == elem.Pos() {
				weCare = true
				break
			}
		}

		if !weCare {
			return m
		}

		// decl doc
		commonCommentText := newDecl("common", elem.Tok)
		if elem.Doc != nil {
			for _, comment := range elem.Doc.List {
				commonCommentText = extractWordAndPhrasesFromComment(commonCommentText, comment.Text, m.dict)
			}
		}

		// names (in case of multiple specs)
		for _, spec := range elem.Specs {
			if valSpec, ok := spec.(*ast.ValueSpec); ok {
				for j, name := range valSpec.Names {
					if name.Name == "_" {
						continue
					}

					declText := newDecl(getDeclID(m.packageName, elem.Tok, name.String()), elem.Tok)
					// TODO add name if valid / split name
					if word := strings.ToLower(name.String()); m.dict.Contains(word) {
						declText.Words[word] = struct{}{}
					}

					if valSpec.Values != nil {
						if val, ok := valSpec.Values[j].(*ast.BasicLit); ok && val.Kind == token.STRING {
							// extract text from value
							valStr := strings.Replace(val.Value, "\"", "", -1)
							// TODO delete commas
							for _, word := range strings.Split(valStr, " ") {
								word = cleaner.ReplaceAllString(word, "")
								if m.dict.Contains(strings.ToLower(word)) {
									declText.Words[strings.ToLower(word)] = struct{}{}
								}
							}

							// get phrases
							phrases, _ := nounphrases.Find(cleanComment(valStr))
							for _, phr := range phrases {
								declText.Phrases[phr] = struct{}{}
							}
						}
					}

					// merge with common comments text
					for k, v := range commonCommentText.Words {
						declText.Words[k] = v
					}

					for k, v := range commonCommentText.Phrases {
						declText.Phrases[k] = v
					}

					m.decls[declText.ID] = declText
				}
			}

			// TODO review spec doc comments!
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				name := typeSpec.Name.String()
				declText := newDecl("", token.TYPE)

				//declText := newText(createFunctionID(m.pkg, token.STRUCT, name), token.TYPE)

				for _, part := range conserv.Split(name) {
					if m.dict.Contains(part) {
						declText.Words[strings.ToLower(part)] = struct{}{}
					}
				}

				// TODO extract field comment
				if typeSpec.Doc != nil {
					for _, comment := range typeSpec.Doc.List {
						declText = extractWordAndPhrasesFromComment(declText, comment.Text, m.dict)
					}
				}

				// TODO extract comments
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					declText.ID = getDeclID(m.packageName, token.STRUCT, name)
					declText.DeclType = token.STRUCT

					if structType.Fields != nil && structType.Fields.List != nil {
						for _, field := range structType.Fields.List {
							for _, fname := range field.Names {
								for _, part := range conserv.Split(fname.Name) {
									if m.dict.Contains(part) {
										declText.Words[strings.ToLower(part)] = struct{}{}
									}
								}
							}

							// TODO extract field comment
							if field.Doc != nil {
								for _, comment := range field.Doc.List {
									declText = extractWordAndPhrasesFromComment(declText, comment.Text, m.dict)
								}
							}
						}
					}
				}

				if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
					declText.ID = getDeclID(m.packageName, token.INTERFACE, name)
					declText.DeclType = token.INTERFACE

					if interfaceType.Methods != nil && interfaceType.Methods.List != nil {
						for _, method := range interfaceType.Methods.List {
							for _, mname := range method.Names {
								for _, part := range conserv.Split(mname.Name) {
									if m.dict.Contains(part) {
										declText.Words[strings.ToLower(part)] = struct{}{}
									}
								}
							}

							// TODO extract method comment
							if method.Doc != nil {
								for _, comment := range method.Doc.List {
									declText = extractWordAndPhrasesFromComment(declText, comment.Text, m.dict)
								}
							}
						}
					}
				}

				// merge with common comments text
				for k, v := range commonCommentText.Words {
					declText.Words[k] = v
				}

				for k, v := range commonCommentText.Phrases {
					declText.Phrases[k] = v
				}

				m.decls[declText.ID] = declText
			}
		}
	}

	return m
}

func getFunctionTextFromFuncDecl(elem *ast.FuncDecl, m Declaration) Decl {
	name := elem.Name.String()
	functionText := newDecl(getDeclID(m.packageName, token.FUNC, name), token.FUNC)

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

func cleanComment(text string) string {
	cleanComment := strings.ReplaceAll(text, "//", "")
	cleanComment = strings.ReplaceAll(cleanComment, "\\n", " ")
	return strings.ReplaceAll(cleanComment, "\\t", " ")
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
		// TODO print log error
	}

	return functionText
}

func getDeclID(pkg string, declType token.Token, name string) string {
	return fmt.Sprintf("%s++%s::%s", pkg, declType, name)
}

// Decls returns a map of declaration IDs and the mined text for each declaration.
func (m Declaration) Decls() map[string]Decl {
	return m.decls
}
