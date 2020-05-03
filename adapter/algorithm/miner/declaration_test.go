package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/adapter/algorithm/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewDeclarationsFactory_ShouldReturnDeclarationsMinerFactory(t *testing.T) {
	factory := miner.NewDeclarationsFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnDeclarationsFactory_ShouldReturnMiner(t *testing.T) {
	factory := miner.NewDeclarationsFactory()
	miner, err := factory.Make()

	assert.Equal(t, "declarations", miner.Name())
	assert.NoError(t, err)
}

func TestSetCurrentFile_OnDeclaration_ShouldSetFile(t *testing.T) {
	factory := miner.NewDeclarationsFactory()
	mnr, _ := factory.Make()

	mnr.SetCurrentFile("new_file.go")

	assert.Equal(t, "new_file.go", mnr.(*miner.Declaration).Filename)
}

func TestVisit_OnDeclarationWithFunctions_ShouldReturnDecls(t *testing.T) {
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
		expected map[string]miner.Decl
	}{
		{"no_functions", "package main", make(map[string]miner.Decl)},
		{"functions_without_text_or_comments", srcWithoutTextOrComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:func+++name:main": {
				ID:       "filename:testfile.go+++pkg:main+++declType:func+++name:main",
				DeclType: token.FUNC,
				Words:    map[string]struct{}{"main": {}},
				Phrases:  make(map[string]struct{}),
			},
		}},
		{"functions_with_text_and_comments", srcWithTextAndComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:func+++name:main": {
				ID:       "filename:testfile.go+++pkg:main+++declType:func+++name:main",
				DeclType: token.FUNC,
				Words: map[string]struct{}{
					"main":     {},
					"is":       {},
					"the":      {},
					"function": {},
					"entry":    {},
					"point":    {},
					"inner":    {},
					"comment":  {},
				},
				Phrases: map[string]struct{}{
					"main function": {},
					"entry point":   {},
					"inner comment": {},
				},
			},
		}},
		{"functions_with_multiple_functions", srcWithMultipleFunctions, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:func+++name:main": {
				ID:       "filename:testfile.go+++pkg:main+++declType:func+++name:main",
				DeclType: token.FUNC,
				Words:    map[string]struct{}{"main": {}},
				Phrases:  make(map[string]struct{}),
			},
			"filename:testfile.go+++pkg:main+++declType:func+++name:another": {
				ID:       "filename:testfile.go+++pkg:main+++declType:func+++name:another",
				DeclType: token.FUNC,
				Words: map[string]struct{}{
					"another":  {},
					"function": {},
					"inside":   {},
				},
				Phrases: make(map[string]struct{}),
			},
		}},
	}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, _ := parser.ParseFile(fs, "testfile.go", []byte(fixture.src), parser.ParseComments)

			factory := miner.NewDeclarationsFactory()
			m, _ := factory.Make()
			m.SetCurrentFile("testfile.go")
			ast.Walk(m, node)

			decls := m.Results().(map[string]miner.Decl)
			assert.Equal(t, len(fixture.expected), len(decls))
			assert.Equal(t, fixture.expected, decls)
		})
	}
}

func TestVisit_OnDeclarationWithVarDecl_ShouldReturnWordsAndPhrases(t *testing.T) {
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
			nrzXXZ int = 32
		)
	`

	tests := []struct {
		name     string
		src      string
		expected map[string]miner.Decl
	}{
		{"variable_without_comments", srcVarWithoutComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:var+++name:Common": {
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:Common",
				DeclType: token.VAR,
				Words:    map[string]struct{}{"common": {}},
				Phrases:  make(map[string]struct{}),
			}}},
		{"variable_with_doc_comments", srcVarWithDocComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:var+++name:Common": {
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:Common",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": {},
					"common":  {},
					"outer":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				},
			}}},
		{"variable_with_doc_and_ignored_line_comments", srcVarWithDocAndLineComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:var+++name:Common": {
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:Common",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": {},
					"common":  {},
					"outer":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				},
			}}},
		{"multiple_variables_same_line", srcMultipleVarSpecs, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:var+++name:common": {
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:common",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": {},
					"common":  {},
					"outer":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				}},
			"filename:testfile.go+++pkg:main+++declType:var+++name:regular": {
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:regular",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": {},
					"regular": {},
					"outer":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				},
			}}},
		{"var_block", srcVarBlock, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:var+++name:common": {
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:common",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": {},
					"common":  {},
					"outer":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				}},
			"filename:testfile.go+++pkg:main+++declType:var+++name:regular": {
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:regular",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": {},
					"regular": {},
					"outer":   {},
					"valid":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				}},
			"filename:testfile.go+++pkg:main+++declType:var+++name:nrzXXZ": {
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:nrzXXZ",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": {},
					"outer":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				},
			}}},
	}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, _ := parser.ParseFile(fs, "", []byte(fixture.src), parser.ParseComments)

			factory := miner.NewDeclarationsFactory()
			m, _ := factory.Make()
			m.SetCurrentFile("testfile.go")
			ast.Walk(m, node)

			decls := m.Results().(map[string]miner.Decl)
			assert.Equal(t, len(fixture.expected), len(decls))
			assert.Equal(t, fixture.expected, decls)
		})
	}
}

func TestVisit_OnDeclarationWithConstDecl_ShouldReturnWordsAndPhrases(t *testing.T) {
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
			nrzXXZ int = 32
		)
	`

	tests := []struct {
		name     string
		src      string
		expected map[string]miner.Decl
	}{
		{"constant_without_comments", srcConstWithoutComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:const+++name:Common": {
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:Common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"common": {},
					"valid":  {},
				},
				Phrases: make(map[string]struct{}),
			}}},
		{"constant_with_doc_comments", srcConstWithDocComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:const+++name:Common": {
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:Common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": {},
					"common":  {},
					"outer":   {},
					"valid":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				},
			}}},
		{"constant_with_doc_and_ignored_line_comments", srcConstWithDocAndLineComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:const+++name:Common": {
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:Common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": {},
					"common":  {},
					"outer":   {},
					"valid":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				},
			}}},
		{"multiple_constants_same_line", srcMultipleConstSpecs, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:const+++name:common": {
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": {},
					"common":  {},
					"outer":   {},
					"valid":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				}},
			"filename:testfile.go+++pkg:main+++declType:const+++name:regular": {
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:regular",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": {},
					"invalid": {},
					"regular": {},
					"outer":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				},
			}}},
		{"const_block", srcConstBlock, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:const+++name:common": {
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": {},
					"common":  {},
					"outer":   {},
					"value":   {},
				},
				Phrases: map[string]struct{}{
					"common value":  {},
					"outer comment": {},
				}},
			"filename:testfile.go+++pkg:main+++declType:const+++name:regular": {
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:regular",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": {},
					"regular": {},
					"outer":   {},
					"valid":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				}},
			"filename:testfile.go+++pkg:main+++declType:const+++name:nrzXXZ": {
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:nrzXXZ",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"comment": {},
					"outer":   {},
				},
				Phrases: map[string]struct{}{
					"outer comment": {},
				},
			}}},
	}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, _ := parser.ParseFile(fs, "", []byte(fixture.src), parser.ParseComments)

			factory := miner.NewDeclarationsFactory()
			m, _ := factory.Make()
			m.SetCurrentFile("testfile.go")
			ast.Walk(m, node)

			decls := m.Results().(map[string]miner.Decl)
			assert.Equal(t, len(fixture.expected), len(decls))
			assert.Equal(t, fixture.expected, decls)
		})
	}
}

func TestVisit_OnDeclarationWithTypeDecl_ShouldReturnWordsAndPhrases(t *testing.T) {
	srcStructWithoutFields := `
		package main

		type selector struct{}
	`

	srcStructWithComments := `
		package main

		// type comment
		type selector struct{}
	`

	srcStructWithFields := `
		package main

		// type comment
		type selector struct {
			pick string
		}
	`

	srcStructWithFieldsAndComments := `
		package main

		// type comment
		type selector struct {
			// field comment
			pick string
		}
	`

	srcStructBlock := `
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

	srcStructWithHardwords := `
		package main

		type httpClient struct {
			protocolPicker string
			url string
		}
	`

	tests := []struct {
		name     string
		src      string
		expected map[string]miner.Decl
	}{
		{"empty_struct", srcStructWithoutFields, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": {
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
				DeclType: token.STRUCT,
				Words: map[string]struct{}{
					"selector": {},
				},
				Phrases: map[string]struct{}{}},
		}},
		{"struct_with_comments", srcStructWithComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": {
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
				DeclType: token.STRUCT,
				Words: map[string]struct{}{
					"comment":  {},
					"selector": {},
					"type":     {},
				},
				Phrases: map[string]struct{}{
					"type comment": {},
				}},
		}},
		{"struct_with_fields", srcStructWithFields, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": {
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
				DeclType: token.STRUCT,
				Words: map[string]struct{}{
					"comment":  {},
					"pick":     {},
					"selector": {},
					"type":     {},
				},
				Phrases: map[string]struct{}{
					"type comment": {},
				}},
		}},
		{"struct_with_fields_and_comments", srcStructWithFieldsAndComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": {
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
				DeclType: token.STRUCT,
				Words: map[string]struct{}{
					"comment":  {},
					"field":    {},
					"pick":     {},
					"selector": {},
					"type":     {},
				},
				Phrases: map[string]struct{}{
					"field comment": {},
					"type comment":  {},
				}},
		}},
		{"struct_block_decl", srcStructBlock, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": {
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
				DeclType: token.STRUCT,
				Words: map[string]struct{}{
					"comment":  {},
					"field":    {},
					"global":   {},
					"local":    {},
					"pick":     {},
					"selector": {},
				},
				Phrases: map[string]struct{}{
					"field comment":  {},
					"global comment": {},
					"local comment":  {},
				}},
			"filename:testfile.go+++pkg:main+++declType:struct+++name:picker": {
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:picker",
				DeclType: token.STRUCT,
				Words: map[string]struct{}{
					"inner":   {},
					"comment": {},
					"global":  {},
					"pick":    {},
					"picker":  {},
				},
				Phrases: map[string]struct{}{
					"inner comment":  {},
					"global comment": {},
				}},
		}},
		{"struct_with_hardwords", srcStructWithHardwords, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:httpClient": {
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:httpClient",
				DeclType: token.STRUCT,
				Words: map[string]struct{}{
					"client":   {},
					"http":     {},
					"picker":   {},
					"protocol": {},
					"url":      {},
				},
				Phrases: map[string]struct{}{}},
		}},
	}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, _ := parser.ParseFile(fs, "", []byte(fixture.src), parser.ParseComments)

			factory := miner.NewDeclarationsFactory()
			m, _ := factory.Make()
			m.SetCurrentFile("testfile.go")
			ast.Walk(m, node)

			decls := m.Results().(map[string]miner.Decl)
			assert.Equal(t, len(fixture.expected), len(decls))
			assert.Equal(t, fixture.expected, decls)
		})
	}
}

func TestVisit_OnDeclarationWithInterfaceDecl_ShouldReturnWordsAndPhrases(t *testing.T) {
	srcInterfaceWithoutMethods := `
		package main

		type selector interface{}
	`

	srcInterfaceWithComments := `
		package main

		// interface comment
		type selector interface{}
	`

	srcInterfaceWithMethods := `
		package main

		// interface comment
		type selector interface {
			pick() string
		}
	`

	srcInterfaceWithMethodsAndComments := `
		package main

		// interface comment
		type selector interface {
			// function comment
			pick() string
		}
	`

	srcInterfaceBlock := `
		package main

		// global comment
		type (
			// interface comment
			selector interface {
				// function comment
				pick() string
			}

			// inner comment
			picker interface {
				pick() string
			}
		)
	`

	srcInterfaceWithHardwords := `
		package main

		type httpClient interface {
			protocolPicker() string
			url() string
		}
	`

	tests := []struct {
		name     string
		src      string
		expected map[string]miner.Decl
	}{
		{"empty_interface", srcInterfaceWithoutMethods, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": {
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"selector": {},
				},
				Phrases: map[string]struct{}{}},
		}},
		{"interface_with_comments", srcInterfaceWithComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": {
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"comment":   {},
					"interface": {},
					"selector":  {},
				},
				Phrases: map[string]struct{}{
					"interface comment": {},
				}},
		}},
		{"interface_with_methods", srcInterfaceWithMethods, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": {
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"comment":   {},
					"interface": {},
					"pick":      {},
					"selector":  {},
				},
				Phrases: map[string]struct{}{
					"interface comment": {},
				}},
		}},
		{"interface_with_methods_and_comments", srcInterfaceWithMethodsAndComments, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": {
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"comment":   {},
					"function":  {},
					"interface": {},
					"pick":      {},
					"selector":  {},
				},
				Phrases: map[string]struct{}{
					"function comment":  {},
					"interface comment": {},
				}},
		}},
		{"interface_block_decl", srcInterfaceBlock, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": {
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"comment":   {},
					"function":  {},
					"global":    {},
					"interface": {},
					"pick":      {},
					"selector":  {},
				},
				Phrases: map[string]struct{}{
					"function comment":  {},
					"global comment":    {},
					"interface comment": {},
				}},
			"filename:testfile.go+++pkg:main+++declType:interface+++name:picker": {
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:picker",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"inner":   {},
					"comment": {},
					"global":  {},
					"pick":    {},
					"picker":  {},
				},
				Phrases: map[string]struct{}{
					"inner comment":  {},
					"global comment": {},
				}},
		}},
		{"interface_with_hardwords", srcInterfaceWithHardwords, map[string]miner.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:httpClient": {
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:httpClient",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"client":   {},
					"http":     {},
					"picker":   {},
					"protocol": {},
					"url":      {},
				},
				Phrases: map[string]struct{}{}},
		}},
	}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, _ := parser.ParseFile(fs, "", []byte(fixture.src), parser.ParseComments)

			factory := miner.NewDeclarationsFactory()
			m, _ := factory.Make()
			m.SetCurrentFile("testfile.go")
			ast.Walk(m, node)

			decls := m.Results().(map[string]miner.Decl)
			assert.Equal(t, len(fixture.expected), len(decls))
			assert.Equal(t, fixture.expected, decls)
		})
	}
}
