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
		{"functions_without_text_or_comments", srcWithoutTextOrComments, map[string]miner.Text{
			"main++func::main": miner.Text{
				ID:       "main++func::main",
				DeclType: token.FUNC,
				Words:    map[string]struct{}{"main": struct{}{}},
				Phrases:  make(map[string]struct{}),
			},
		}},
		{"functions_with_text_and_comments", srcWithTextAndComments, map[string]miner.Text{
			"main++func::main": miner.Text{
				ID:       "main++func::main",
				DeclType: token.FUNC,
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
		{"functions_with_multiple_functions", srcWithMultipleFunctions, map[string]miner.Text{
			"main++func::main": miner.Text{
				ID:       "main++func::main",
				DeclType: token.FUNC,
				Words:    map[string]struct{}{"main": struct{}{}},
				Phrases:  make(map[string]struct{}),
			},
			"main++func::another": miner.Text{
				ID:       "main++func::another",
				DeclType: token.FUNC,
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

func TestVisit_OnFunctionWithVarDecl_ShouldReturnWordsAndPhrases(t *testing.T) {
	srcVarWithoutComments := `
		package main
		
		var Common string
	`

	srcVarWithDocComments := `
		package main

		// outer comment
		var Common string
	`
	srcVarWithDocAndLineComments := `
		package main

		// outer comment
		var Common string // inner comment
	`

	srcMultipleVarSpecs := `
		package main

		// outer comment
		var common,regular string
	`

	srcVarBlock := `
		package main

		// outer comment
		var (
			common string
			regular string = "valid"
			nrXX int = 32
		)
	`

	tests := []struct {
		name     string
		src      string
		expected map[string]miner.Text
	}{
		{"variable_without_comments", srcVarWithoutComments, map[string]miner.Text{
			"main++var::Common": miner.Text{
				ID:       "main++var::Common",
				DeclType: token.VAR,
				Words:    map[string]struct{}{"common": struct{}{}},
				Phrases:  make(map[string]struct{}),
			}}},
		{"variable_with_doc_comments", srcVarWithDocComments, map[string]miner.Text{
			"main++var::Common": miner.Text{
				ID:       "main++var::Common",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"common":  struct{}{},
					"outer":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				},
			}}},
		{"variable_with_doc_and_ignored_line_comments", srcVarWithDocAndLineComments, map[string]miner.Text{
			"main++var::Common": miner.Text{
				ID:       "main++var::Common",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"common":  struct{}{},
					"outer":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				},
			}}},
		{"multiple_variables_same_line", srcMultipleVarSpecs, map[string]miner.Text{
			"main++var::common": miner.Text{
				ID:       "main++var::common",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"common":  struct{}{},
					"outer":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				}},
			"main++var::regular": miner.Text{
				ID:       "main++var::regular",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"regular": struct{}{},
					"outer":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				},
			}}},
		{"var_block", srcVarBlock, map[string]miner.Text{
			"main++var::common": miner.Text{
				ID:       "main++var::common",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"common":  struct{}{},
					"outer":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				}},
			"main++var::regular": miner.Text{
				ID:       "main++var::regular",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"regular": struct{}{},
					"outer":   struct{}{},
					"valid":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				}},
			"main++var::nrXX": miner.Text{
				ID:       "main++var::nrXX",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"outer":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				},
			}}},
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

/*func TestVisit_OnFunctionWithRealLifeFile_ShouldReturnFunctionsWordsAndPhrases(t *testing.T) {
	src := `
	`
	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "", []byte(src), parser.ParseComments)

	function := miner.NewFunction(lists.Dicctionary)
	ast.Walk(function, node)

	functions := function.FunctionsText()
	assert.Equal(t, 3, len(functions))
	// assert.Equal(t, fixture.expected, functions)
}*/
