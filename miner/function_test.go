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

func TestVisit_OnFunctionWithConstDecl_ShouldReturnWordsAndPhrases(t *testing.T) {
	srcConstWithoutComments := `
		package main

		const Common string = "valid"
	`

	srcConstWithDocComments := `
		package main

		// outer comment
		const Common string = "valid"
	`
	srcConstWithDocAndLineComments := `
		package main

		// outer comment
		const Common string = "valid" // inner comment
	`

	srcMultipleConstSpecs := `
		package main

		// outer comment
		const common,regular string = "valid", "invalid"
	`

	srcConstBlock := `
		package main

		// outer comment
		const (
			common string = "common value"
			regular string = "valid"
			nrXX int = 32
		)
	`

	tests := []struct {
		name     string
		src      string
		expected map[string]miner.Text
	}{
		{"constant_without_comments", srcConstWithoutComments, map[string]miner.Text{
			"main++const::Common": miner.Text{
				ID:       "main++const::Common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"common": struct{}{},
					"valid":  struct{}{},
				},
				Phrases: make(map[string]struct{}),
			}}},
		{"constant_with_doc_comments", srcConstWithDocComments, map[string]miner.Text{
			"main++const::Common": miner.Text{
				ID:       "main++const::Common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"common":  struct{}{},
					"outer":   struct{}{},
					"valid":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				},
			}}},
		{"constant_with_doc_and_ignored_line_comments", srcConstWithDocAndLineComments, map[string]miner.Text{
			"main++const::Common": miner.Text{
				ID:       "main++const::Common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"common":  struct{}{},
					"outer":   struct{}{},
					"valid":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				},
			}}},
		{"multiple_constants_same_line", srcMultipleConstSpecs, map[string]miner.Text{
			"main++const::common": miner.Text{
				ID:       "main++const::common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"common":  struct{}{},
					"outer":   struct{}{},
					"valid":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				}},
			"main++const::regular": miner.Text{
				ID:       "main++const::regular",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"invalid": struct{}{},
					"regular": struct{}{},
					"outer":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				},
			}}},
		{"const_block", srcConstBlock, map[string]miner.Text{
			"main++const::common": miner.Text{
				ID:       "main++const::common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"common":  struct{}{},
					"outer":   struct{}{},
					"value":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"common value":  struct{}{},
					"outer comment": struct{}{},
				}},
			"main++const::regular": miner.Text{
				ID:       "main++const::regular",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"regular": struct{}{},
					"outer":   struct{}{},
					"valid":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				}},
			"main++const::nrXX": miner.Text{
				ID:       "main++const::nrXX",
				DeclType: token.CONST,
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

func TestVisit_OnFunctionWithTypeDecl_ShouldReturnWordsAndPhrases(t *testing.T) {
	srcTypeWithoutFields := `
		package main

		type selector struct{}
	`

	srcTypeWithComments := `
		package main

		// type comment
		type selector struct{}
	`

	srcTypeWithFields := `
		package main

		// type comment
		type selector struct {
			pick string
		}
	`

	srcTypeWithFieldsAndComments := `
		package main

		// type comment
		type selector struct {
			// field comment
			pick string
		}
	`

	srcTypeBlock := `
		package main

		// global comment
		type (
			// local comment
			selector struct {
				// field comment
				pick string
			}

			// inner comment
			picker struct {
				pick string
			}
		)
	`

	srcTypeWithHardwords := `
		package main

		type httpClient struct {
			protocolPicker string
			url string
		}
	`

	tests := []struct {
		name     string
		src      string
		expected map[string]miner.Text
	}{
		{"empty_type", srcTypeWithoutFields, map[string]miner.Text{
			"main++type::selector": miner.Text{
				ID:       "main++type::selector",
				DeclType: token.TYPE,
				Words: map[string]struct{}{
					"selector": struct{}{},
				},
				Phrases: map[string]struct{}{}},
		}},
		{"type_with_comments", srcTypeWithComments, map[string]miner.Text{
			"main++type::selector": miner.Text{
				ID:       "main++type::selector",
				DeclType: token.TYPE,
				Words: map[string]struct{}{
					"comment":  struct{}{},
					"selector": struct{}{},
					"type":     struct{}{},
				},
				Phrases: map[string]struct{}{
					"type comment": struct{}{},
				}},
		}},
		{"type_with_fields", srcTypeWithFields, map[string]miner.Text{
			"main++type::selector": miner.Text{
				ID:       "main++type::selector",
				DeclType: token.TYPE,
				Words: map[string]struct{}{
					"comment":  struct{}{},
					"pick":     struct{}{},
					"selector": struct{}{},
					"type":     struct{}{},
				},
				Phrases: map[string]struct{}{
					"type comment": struct{}{},
				}},
		}},
		{"type_with_fields_and_comments", srcTypeWithFieldsAndComments, map[string]miner.Text{
			"main++type::selector": miner.Text{
				ID:       "main++type::selector",
				DeclType: token.TYPE,
				Words: map[string]struct{}{
					"comment":  struct{}{},
					"field":    struct{}{},
					"pick":     struct{}{},
					"selector": struct{}{},
					"type":     struct{}{},
				},
				Phrases: map[string]struct{}{
					"field comment": struct{}{},
					"type comment":  struct{}{},
				}},
		}},
		{"type_block_decl", srcTypeBlock, map[string]miner.Text{
			"main++type::selector": miner.Text{
				ID:       "main++type::selector",
				DeclType: token.TYPE,
				Words: map[string]struct{}{
					"comment":  struct{}{},
					"field":    struct{}{},
					"global":   struct{}{},
					"local":    struct{}{},
					"pick":     struct{}{},
					"selector": struct{}{},
				},
				Phrases: map[string]struct{}{
					"field comment":  struct{}{},
					"global comment": struct{}{},
					"local comment":  struct{}{},
				}},
			"main++type::picker": miner.Text{
				ID:       "main++type::picker",
				DeclType: token.TYPE,
				Words: map[string]struct{}{
					"inner":   struct{}{},
					"comment": struct{}{},
					"global":  struct{}{},
					"pick":    struct{}{},
					"picker":  struct{}{},
				},
				Phrases: map[string]struct{}{
					"inner comment":  struct{}{},
					"global comment": struct{}{},
				}},
		}},
		{"type_with_hardwords", srcTypeWithHardwords, map[string]miner.Text{
			"main++type::httpClient": miner.Text{
				ID:       "main++type::httpClient",
				DeclType: token.TYPE,
				Words: map[string]struct{}{
					"client":   struct{}{},
					"http":     struct{}{},
					"picker":   struct{}{},
					"protocol": struct{}{},
					"url":      struct{}{},
				},
				Phrases: map[string]struct{}{}},
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
