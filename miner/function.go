package miner

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/eroatta/nounphrases"
	"github.com/eroatta/token/lists"
)

// Text contains the mined text extracted from a function.
type Text struct {
	ID       string
	DeclType token.Token
	Words    map[string]struct{}
	Phrases  map[string]struct{}
}

func newText(ID string, declType token.Token) Text {
	return Text{
		ID:       ID,
		DeclType: declType,
		Words:    make(map[string]struct{}),
		Phrases:  make(map[string]struct{}),
	}
}

// Function represents the functions miner, which extracts information about
// words and phrases for each function declaration.
type Function struct {
	dict        lists.List
	pkg         string
	pkgComments []*ast.CommentGroup
	decls       []ast.Decl
	functions   map[string]Text
}

// NewFunction initializes a new functions miner.
func NewFunction(dict lists.List) Function {
	return Function{
		dict:      dict,
		functions: make(map[string]Text),
	}
}

// Name returns the specific name for the miner.
func (m Function) Name() string {
	return "function"
}

// Visit implements the ast.Visitor interface and handles the logic for the data extraction.
func (m Function) Visit(node ast.Node) ast.Visitor {
	// on a FuncDecl, get comments + identifiers + body text
	// and extract words (stemming + stopping)
	// and also extract phrases
	if node == nil {
		return nil
	}

	switch elem := node.(type) {
	case *ast.File:
		m.pkg = elem.Name.String()
		m.decls = elem.Decls

		// collect inner comment to use on their proper function
		m.pkgComments = append(m.pkgComments, elem.Comments...)

	case *ast.FuncDecl:
		name := elem.Name.String()
		functionText := newText(createFunctionID(m.pkg, token.FUNC, name), token.FUNC)

		if m.dict.Contains(name) {
			functionText.Words[strings.ToLower(name)] = struct{}{}
		}

		if elem.Doc != nil {
			for _, comment := range elem.Doc.List {
				functionText = extractWordAndPhrasesFromComment(functionText, comment.Text, m.dict)
			}
		}

		start, end := elem.Pos(), elem.End()
		for _, group := range m.pkgComments {
			for _, comment := range group.List {
				if comment.Slash > start && comment.Slash < end {
					functionText = extractWordAndPhrasesFromComment(functionText, comment.Text, m.dict)
				}
			}
		}

		m.functions[functionText.ID] = functionText

	case *ast.GenDecl:
		// TODO improve this code
		if elem.Tok != token.VAR {
			return m
		}

		var weCare bool
		for _, d := range m.decls {
			if d.Pos() == elem.Pos() {
				weCare = true
				break
			}
		}

		if !weCare {
			return m
		}

		// names (in case of multiple specs)
		var decls []Text
		var start int // TODO review
		for _, spec := range elem.Specs {
			if valSpec, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range valSpec.Names {
					if name.Name == "_" {
						continue
					}

					declText := newText(createFunctionID(m.pkg, token.VAR, name.String()), token.VAR)
					// TODO add name if valid
					if word := strings.ToLower(name.String()); m.dict.Contains(word) {
						declText.Words[word] = struct{}{}
					}

					decls = append(decls, declText)
					start++
				}

				fmt.Println(start)
				for i, value := range valSpec.Values {
					if val, ok := value.(*ast.BasicLit); ok && val.Kind == token.STRING {
						// extract text from value
						valStr := strings.Replace(val.Value, "\"", "", -1)
						// TODO delete commas
						for _, word := range strings.Split(valStr, " ") {
							word = cleaner.ReplaceAllString(word, "")
							if m.dict.Contains(strings.ToLower(word)) {
								fmt.Println(fmt.Sprintf("Start %d, i %d, word: %s", start, i, word))
								fmt.Println(len(decls))
								decls[start+i-1].Words[strings.ToLower(word)] = struct{}{}
							}
						}

						// get phrases
						phrases, _ := nounphrases.Find(cleanComment(valStr))
						for _, phr := range phrases {
							decls[start+i-1].Phrases[phr] = struct{}{} // TODO: review index
						}
					}
				}
			}
		}

		// decl doc
		dummyCommentText := newText("", token.VAR)
		if elem.Doc != nil {
			for _, comment := range elem.Doc.List {
				dummyCommentText = extractWordAndPhrasesFromComment(dummyCommentText, comment.Text, m.dict)
			}
		}

		// merge decl doc
		for i := range decls {
			for k, v := range dummyCommentText.Words {
				decls[i].Words[k] = v
			}

			for k, v := range dummyCommentText.Phrases {
				decls[i].Phrases[k] = v
			}

			// add them to the map
			m.functions[decls[i].ID] = decls[i]
		}
	}

	return m
}

func cleanComment(text string) string {
	cleanComment := strings.ReplaceAll(text, "//", "")
	cleanComment = strings.ReplaceAll(cleanComment, "\\n", " ")
	return strings.ReplaceAll(cleanComment, "\\t", " ")
}

func extractWordAndPhrasesFromComment(functionText Text, comment string, list lists.List) Text {
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

func createFunctionID(pkg string, declType token.Token, name string) string {
	return fmt.Sprintf("%s++%s::%s", pkg, declType, name)
}

// FunctionsText returns a map of function names and the mined text for each function.
func (m Function) FunctionsText() map[string]Text {
	return m.functions
}
