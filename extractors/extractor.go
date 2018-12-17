package extractors

import "go/ast"

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
	name string
}

// NewSamuraiExtractor creates an instance capable of exploring the Abstract Systax Tree
// and extracting the data related to the Samurai splitting algorithm.
func NewSamuraiExtractor() Extractor {
	return SamuraiExtractor{name: "samurai"}
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

	return e
}
