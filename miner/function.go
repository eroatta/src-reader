package miner

import "go/ast"

// Text contains the mined text extracted from a function.
type Text struct {
	Words   map[string]struct{}
	Phrases map[string]string
}

// Function TODO
type Function struct {
	funcs map[string]Text
}

// NewFunction TODO
func NewFunction() Function {
	return Function{
		funcs: make(map[string]Text),
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
	return m
}

// FunctionsText returns a map of function names and the mined text for each function.
func (m Function) FunctionsText() map[string]Text {
	return m.funcs
}
