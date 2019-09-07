package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/token/lists"
	"github.com/stretchr/testify/assert"
)

func TestNewFunction_ShouldReturnFunctionMiner(t *testing.T) {
	miner := miner.NewFunction(nil)

	assert.NotNil(t, miner)
}

func TestGetName_OnFunction_ShouldReturnFunction(t *testing.T) {
	miner := miner.NewFunction(nil)

	assert.Equal(t, "function", miner.Name())
}

func TestVisit_OnFunction_ShouldReturnFunctionsWordsAndPhrases(t *testing.T) {
	srcWithoutTextOrComments := `
		package main

		func main() {

		}
	`

	srcWithTextAndComments := `
		package main

		// main is the main function, the entry point, zarasa
		func main() {
			// inner comment
		}
	`

	srcWithMultipleFunctions := `
		package main

		func main() {

		}

		func another() {
			// inside another function
		}
	`

	tests := []struct {
		name     string
		src      string
		expected map[string]miner.Text
	}{
		{"no_functions", "package main", make(map[string]miner.Text)},
		{"without_text_or_comments", srcWithoutTextOrComments, map[string]miner.Text{
			"main+main": miner.Text{
				ID:      "main+main",
				Words:   map[string]struct{}{"main": struct{}{}},
				Phrases: make(map[string]struct{}),
			},
		}},
		{"with_text_and_comments", srcWithTextAndComments, map[string]miner.Text{
			"main+main": miner.Text{
				ID: "main+main",
				Words: map[string]struct{}{
					"main":     struct{}{},
					"is":       struct{}{},
					"the":      struct{}{},
					"function": struct{}{},
					"entry":    struct{}{},
					"point":    struct{}{},
					"inner":    struct{}{},
					"comment":  struct{}{},
				},
				Phrases: map[string]struct{}{
					"main function": struct{}{},
					"entry point":   struct{}{},
					"inner comment": struct{}{},
				},
			},
		}},
		{"with_multiple_functions", srcWithMultipleFunctions, map[string]miner.Text{
			"main+main": miner.Text{
				ID:      "main+main",
				Words:   map[string]struct{}{"main": struct{}{}},
				Phrases: make(map[string]struct{}),
			},
			"main+another": miner.Text{
				ID: "main+another",
				Words: map[string]struct{}{
					"another":  struct{}{},
					"function": struct{}{},
					"inside":   struct{}{},
				},
				Phrases: make(map[string]struct{}),
			},
		}},
	}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, _ := parser.ParseFile(fs, "", []byte(fixture.src), parser.ParseComments)

			function := miner.NewFunction(lists.Dicctionary)
			ast.Walk(function, node)

			functions := function.FunctionsText()
			assert.Equal(t, len(fixture.expected), len(functions))
			assert.Equal(t, fixture.expected, functions)
		})
	}
}

func TestVisit_OnFunctionWithRealLifeFile_ShouldReturnFunctionsWordsAndPhrases(t *testing.T) {
	src := `
	`
	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "", []byte(src), parser.ParseComments)

	function := miner.NewFunction(lists.Dicctionary)
	ast.Walk(function, node)

	functions := function.FunctionsText()
	assert.Equal(t, 3, len(functions))
	// assert.Equal(t, fixture.expected, functions)
}
