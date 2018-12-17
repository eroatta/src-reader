package visitors

import "go/ast"

// Visitor represents a set of ast.Visitors that can be applied to a ast.Node.
type Visitor struct {
	Processors []ast.Visitor
}

func NewVisitor(processors []ast.Visitor) ast.Visitor {
	return Visitor{Processors: processors}
}

// Visit implements the ast.Visitor interface and handles the execution of all the visitors attached to it.
func (v Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	for _, visitor := range v.Processors {
		visitor.Visit(node)
	}

	return v
}
