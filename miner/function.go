package miner

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/eroatta/nounphrases"
	"github.com/eroatta/token/lists"
)

// Text contains the mined text extracted from a function.
type Text struct {
	ID      string
	Words   map[string]struct{}
	Phrases map[string]struct{}
}

func newText(ID string) Text {
	return Text{
		ID:      ID,
		Words:   make(map[string]struct{}),
		Phrases: make(map[string]struct{}),
	}
}

// Function represents the functions miner, which extracts information about
// words and phrases for each function declaration.
type Function struct {
	dict        lists.List
	pkg         string
	pkgComments []*ast.CommentGroup
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

		// collect inner comment to use on their proper function
		m.pkgComments = append(m.pkgComments, elem.Comments...)

	case *ast.FuncDecl:
		name := elem.Name.String()
		functionText := newText(createFunctionID(m.pkg, name))

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

func createFunctionID(pkg string, name string) string {
	return fmt.Sprintf("%s+%s", pkg, name)
}

// FunctionsText returns a map of function names and the mined text for each function.
func (m Function) FunctionsText() map[string]Text {
	return m.functions
}
