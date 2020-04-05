package miner

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/conserv"
)

var cleaner = regexp.MustCompile("[^a-zA-Z0-9]")

// WordCount handles the word count mining process.
type WordCount struct {
	words map[string]int
}

// NewWordCount creates a new Count miner.
func NewWordCount() WordCount {
	return WordCount{
		words: map[string]int{},
	}
}

// Type returns the miner type.
func (m WordCount) Type() entity.MinerType {
	return entity.WordCount
}

// Visit implements the ast.Visitor interface and handles the logic for the data extraction.
func (m WordCount) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	var tokens []string

	switch elem := node.(type) {
	case *ast.AssignStmt:
		tokens = append(tokens, countOnAssignment(elem)...)

	case *ast.RangeStmt:
		tokens = append(tokens, countOnRange(elem)...)

	case *ast.GenDecl:
		switch elem.Tok {
		case token.VAR, token.CONST:
			tokens = append(tokens, countOnVarConstDecl(elem)...)
		case token.TYPE:
			tokens = append(tokens, countOnTypeDecl(elem)...)
		default:
			return m
		}

	case *ast.FuncDecl:
		tokens = append(tokens, countOnFuncDecl(elem)...)

	case *ast.File:
		tokens = append(tokens, countOnFile(elem)...)
	}

	for _, token := range tokens {
		for _, splitting := range conserv.Split(token) {
			m.words[strings.ToLower(splitting)]++
		}
	}

	return m
}

func countOnAssignment(elem *ast.AssignStmt) []string {
	if elem.Tok != token.DEFINE {
		return []string{}
	}

	tokens := []string{}
	for _, expr := range elem.Lhs {
		if identifier, ok := expr.(*ast.Ident); ok {
			if identifier.String() == "_" {
				continue
			}

			// only newly defined identifiers
			if identifier.Obj != nil && identifier.Obj.Pos() == identifier.Pos() {
				tokens = append(tokens, identifier.String())
			}
		}
	}

	return tokens
}

func countOnRange(elem *ast.RangeStmt) []string {
	tokens := []string{}
	if key, ok := elem.Key.(*ast.Ident); ok {
		if key.String() != "_" {
			tokens = append(tokens, key.String())
		}
	}

	if value, ok := elem.Value.(*ast.Ident); ok {
		if value.String() != "_" {
			tokens = append(tokens, value.String())
		}
	}

	return tokens
}

func countOnVarConstDecl(elem *ast.GenDecl) []string {
	tokens := []string{}
	for _, spec := range elem.Specs {
		if valSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range valSpec.Names {
				if name.Name == "_" {
					continue
				}

				tokens = append(tokens, name.String())
			}

			for _, value := range valSpec.Values {
				if val, ok := value.(*ast.BasicLit); ok && val.Kind == token.STRING {
					tokens = append(tokens, strings.Replace(val.Value, "\"", "", -1))
				}
			}
		}
	}

	return tokens
}

func countOnTypeDecl(elem *ast.GenDecl) []string {
	tokens := []string{}
	for _, spec := range elem.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			tokens = append(tokens, typeSpec.Name.String())

			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				for _, field := range structType.Fields.List {
					for _, fieldName := range field.Names {
						tokens = append(tokens, fieldName.String())
					}
				}
			}

			if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				for _, method := range interfaceType.Methods.List {
					for _, name := range method.Names {
						tokens = append(tokens, name.String())
					}

					if funcType, ok := method.Type.(*ast.FuncType); ok {
						for _, in := range funcType.Params.List {
							for _, arg := range in.Names {
								tokens = append(tokens, arg.String())
							}
						}

						results := funcType.Results
						if results != nil {
							for _, out := range results.List {
								for _, arg := range out.Names {
									tokens = append(tokens, arg.String())
								}
							}
						}
					}
				}
			}
		}
	}

	return tokens
}

func countOnFuncDecl(elem *ast.FuncDecl) []string {
	tokens := []string{elem.Name.String()}
	for _, in := range elem.Type.Params.List {
		for _, arg := range in.Names {
			tokens = append(tokens, arg.String())
		}
	}

	results := elem.Type.Results
	if results != nil {
		for _, out := range results.List {
			for _, arg := range out.Names {
				tokens = append(tokens, arg.String())
			}
		}
	}

	return tokens
}

func countOnFile(elem *ast.File) []string {
	tokens := []string{}
	for _, commentGroup := range elem.Comments {
		for _, comment := range commentGroup.List {
			cleanComment := strings.Trim(cleaner.ReplaceAllString(comment.Text, " "), "")
			for _, word := range strings.Split(cleanComment, " ") {
				if word == "" {
					continue
				}

				tokens = append(tokens, word)
			}
		}
	}

	return tokens
}

// Results returns the word count.
func (m WordCount) Results() map[string]int {
	return m.words
}
