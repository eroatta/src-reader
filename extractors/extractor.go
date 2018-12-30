package extractors

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"

	"github.com/eroatta/token-splitex/splitters"
)

// Process traverses the Abstract Systax Tree node and applies the extraction method defined by the extractor.
func Process(extractor Extractor, node ast.Node) {
	extractor.Visit(node)
}

// Extractor defines TODO
type Extractor interface {
	ast.Visitor
	Name() string // name of the extractor.
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

	// TODO: remove after development
	log.Println(node)

	switch elem := node.(type) {
	case *ast.GenDecl:
		if elem.Tok == token.VAR {
			for _, spec := range elem.Specs {
				if valSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valSpec.Names {
						if name.Name == "_" {
							continue
						}

						tokens = append(tokens, name.Name)
					}
				}
			}
		}

		if elem.Tok == token.CONST {
			for _, spec := range elem.Specs {
				if valSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valSpec.Names {
						if name.Name == "_" {
							continue
						}

						tokens = append(tokens, name.Name)
					}

					for _, value := range valSpec.Values {
						if val, ok := value.(*ast.BasicLit); ok && val.Kind == token.STRING {
							tokens = append(tokens, strings.Replace(val.Value, "\"", "", -1))
						}
					}
				}
			}
		}

		// for _, spec := range elem.Specs {
		// 	if valSpec, ok := spec.(*ast.ValueSpec); ok {
		// 		for _, name := range valSpec.Names {
		// 			if name.Name == "_" {
		// 				continue
		// 			}

		// 			tokens = append(tokens, name.Name)
		// 		}
		// 	}
		// }
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
