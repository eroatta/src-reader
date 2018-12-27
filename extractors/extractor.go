package extractors

import (
	"go/ast"
	"log"
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
	name  string
	words map[string]int
}

// NewSamuraiExtractor creates an instance capable of exploring the Abstract Systax Tree
// and extracting the data related to the Samurai splitting algorithm.
func NewSamuraiExtractor() Extractor {
	return SamuraiExtractor{name: "samurai", words: map[string]int{}}
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

	log.Println(node)
	switch elem := node.(type) {
	case *ast.GenDecl:
		decl := *elem
		//log.Println(decl.Tok.String())
		if len(decl.Specs) > 0 {
			spec, _ := decl.Specs[0].(*ast.ValueSpec)
			e.words[spec.Names[0].Name]++
		}
	case *ast.Ident:
		ident := *elem
		e.words[ident.Name] = 1
	}

	return e
}
