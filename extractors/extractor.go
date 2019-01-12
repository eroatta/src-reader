package extractors

import "go/ast"

// Process traverses the Abstract Systax Tree node and applies the extraction method defined by the extractor.
func Process(extractor Extractor, node ast.Node) {
	ast.Walk(extractor, node)
}

// Extractor defines TODO
type Extractor interface {
	ast.Visitor
	Name() string // name of the extractor.
}
