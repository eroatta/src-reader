package extractors

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"regexp"
	"strings"

	"github.com/eroatta/token-splitex/splitters"
)

var cleaner *regexp.Regexp

func init() {
	cleaner = regexp.MustCompile("[^a-zA-Z0-9]")
}

// SamuraiExtractor represents an extractor that reads and stores the required data for the Samurai
// splitting algorithm.
type SamuraiExtractor struct {
	name     string
	words    map[string]int
	splitter splitters.Splitter
}

// NewSamuraiExtractor creates an instance capable of exploring the Abstract Systax Tree
// and extracting the data related to the Samurai splitting algorithm.
func NewSamuraiExtractor() Extractor {
	return SamuraiExtractor{
		name:     "samurai",
		words:    map[string]int{},
		splitter: splitters.NewConserv(),
	}
}

// Name returns the specific name for the extractor.
func (e SamuraiExtractor) Name() string {
	return e.name
}

// Visit implements the ast.Visitor interface and handles the logic for the data extraction.
func (e SamuraiExtractor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	var tokens []string

	switch elem := node.(type) {
	case *ast.AssignStmt:
		if elem.Tok != token.DEFINE {
			return e
		}

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

	case *ast.RangeStmt:
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

	case *ast.GenDecl:
		switch elem.Tok {
		case token.VAR, token.CONST:
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
		case token.TYPE:
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
		default:
			return e
		}

	case *ast.FuncDecl:
		tokens = append(tokens, elem.Name.String())

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

	case *ast.File:
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
	}

	for _, token := range tokens {
		splittings, err := e.splitter.Split(token)
		if err != nil {
			log.Println(fmt.Sprintf("An error occurred while splitting %s: %v", token, err))
			continue
		}

		for _, splitting := range splittings {
			e.words[strings.ToLower(splitting)]++
		}
	}

	return e
}
