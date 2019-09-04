package miner

import "go/ast"

// FunctionData contains the mined function data.
type FunctionData struct {
	ID      string
	words   map[string]struct{}
	phrases map[string]string
}

// Function TODO
type Function struct {
	funcs map[string]FunctionData
}

// NewFunction TODO
func NewFunction() Function {
	return Function{
		funcs: make(map[string]FunctionData),
	}
}

// Name returns the specific name for the miner.
func (m Function) Name() string {
	return "function"
}

// Visit implements the ast.Visitor interface and handles the logic for the data extraction.
func (m Function) Visit(node ast.Node) ast.Visitor {
	return m
}
