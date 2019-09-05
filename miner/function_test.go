package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewFunction_ShouldReturnFunctionMiner(t *testing.T) {
	miner := miner.NewFunction()

	assert.NotNil(t, miner)
}

func TestGetName_OnFunction_ShouldReturnFunction(t *testing.T) {
	miner := miner.NewFunction()

	assert.Equal(t, "function", miner.Name())
}

func TestVisit_OnFunction_ShouldReturnFunctionsWordsAndPhrases(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected map[string]miner.Text
	}{
		{"no_functions", "package main", make(map[string]miner.Text)},
		// TODO complete with other test cases
	}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, _ := parser.ParseFile(fs, "", []byte(fixture.src), parser.ParseComments)

			function := miner.NewFunction()
			ast.Walk(function, node)

			functions := function.FunctionsText()
			assert.Equal(t, len(fixture.expected), len(functions))
			assert.Equal(t, fixture.expected, functions)
		})
	}
}
