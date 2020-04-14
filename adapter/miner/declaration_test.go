package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/adapter/miner"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/lists"
	"github.com/stretchr/testify/assert"
)

func TestNewDeclaration_ShouldReturnDeclarationMiner(t *testing.T) {
	miner := miner.NewDeclaration(nil)

	assert.NotNil(t, miner)
}

func TestType_OnDeclaration_ShouldReturnDeclaration(t *testing.T) {
	miner := miner.NewDeclaration(nil)

	assert.Equal(t, entity.MinerDeclarations, miner.Type())
}

func TestSetCurrentFile_OnDeclaration_ShouldSetFile(t *testing.T) {
	miner := miner.NewDeclaration(nil)
	miner.SetCurrentFile("new_file.go")

	assert.Equal(t, "new_file.go", miner.Filename)
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
		expected map[string]entity.Decl
	}{
		{"no_functions", "package main", make(map[string]entity.Decl)},
		{"functions_without_text_or_comments", srcWithoutTextOrComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:func+++name:main": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:func+++name:main",
				DeclType: token.FUNC,
				Words:    map[string]struct{}{"main": struct{}{}},
				Phrases:  make(map[string]struct{}),
			},
		}},
		{"functions_with_text_and_comments", srcWithTextAndComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:func+++name:main": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:func+++name:main",
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
		{"functions_with_multiple_functions", srcWithMultipleFunctions, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:func+++name:main": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:func+++name:main",
				DeclType: token.FUNC,
				Words:    map[string]struct{}{"main": struct{}{}},
				Phrases:  make(map[string]struct{}),
			},
			"filename:testfile.go+++pkg:main+++declType:func+++name:another": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:func+++name:another",
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
			node, _ := parser.ParseFile(fs, "testfile.go", []byte(fixture.src), parser.ParseComments)

			m := miner.NewDeclaration(lists.Dictionary)
			m.SetCurrentFile("testfile.go")
			ast.Walk(m, node)

			decls := m.Declarations()
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
		expected map[string]entity.Decl
	}{
		{"variable_without_comments", srcVarWithoutComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:var+++name:Common": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:Common",
				DeclType: token.VAR,
				Words:    map[string]struct{}{"common": struct{}{}},
				Phrases:  make(map[string]struct{}),
			}}},
		{"variable_with_doc_comments", srcVarWithDocComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:var+++name:Common": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:Common",
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
		{"variable_with_doc_and_ignored_line_comments", srcVarWithDocAndLineComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:var+++name:Common": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:Common",
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
		{"multiple_variables_same_line", srcMultipleVarSpecs, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:var+++name:common": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:common",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"common":  struct{}{},
					"outer":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				}},
			"filename:testfile.go+++pkg:main+++declType:var+++name:regular": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:regular",
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
		{"var_block", srcVarBlock, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:var+++name:common": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:common",
				DeclType: token.VAR,
				Words: map[string]struct{}{
					"comment": struct{}{},
					"common":  struct{}{},
					"outer":   struct{}{},
				},
				Phrases: map[string]struct{}{
					"outer comment": struct{}{},
				}},
			"filename:testfile.go+++pkg:main+++declType:var+++name:regular": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:regular",
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
			"filename:testfile.go+++pkg:main+++declType:var+++name:nrzXXZ": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:var+++name:nrzXXZ",
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

			m := miner.NewDeclaration(lists.Dictionary)
			m.SetCurrentFile("testfile.go")
			ast.Walk(m, node)

			decls := m.Declarations()
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
		expected map[string]entity.Decl
	}{
		{"constant_without_comments", srcConstWithoutComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:const+++name:Common": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:Common",
				DeclType: token.CONST,
				Words: map[string]struct{}{
					"common": struct{}{},
					"valid":  struct{}{},
				},
				Phrases: make(map[string]struct{}),
			}}},
		{"constant_with_doc_comments", srcConstWithDocComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:const+++name:Common": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:Common",
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
		{"constant_with_doc_and_ignored_line_comments", srcConstWithDocAndLineComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:const+++name:Common": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:Common",
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
		{"multiple_constants_same_line", srcMultipleConstSpecs, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:const+++name:common": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:common",
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
			"filename:testfile.go+++pkg:main+++declType:const+++name:regular": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:regular",
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
		{"const_block", srcConstBlock, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:const+++name:common": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:common",
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
			"filename:testfile.go+++pkg:main+++declType:const+++name:regular": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:regular",
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
			"filename:testfile.go+++pkg:main+++declType:const+++name:nrzXXZ": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:const+++name:nrzXXZ",
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

			m := miner.NewDeclaration(lists.Dictionary)
			m.SetCurrentFile("testfile.go")
			ast.Walk(m, node)

			decls := m.Declarations()
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
		expected map[string]entity.Decl
	}{
		{"empty_struct", srcStructWithoutFields, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
				DeclType: token.STRUCT,
				Words: map[string]struct{}{
					"selector": struct{}{},
				},
				Phrases: map[string]struct{}{}},
		}},
		{"struct_with_comments", srcStructWithComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
				DeclType: token.STRUCT,
				Words: map[string]struct{}{
					"comment":  struct{}{},
					"selector": struct{}{},
					"type":     struct{}{},
				},
				Phrases: map[string]struct{}{
					"type comment": struct{}{},
				}},
		}},
		{"struct_with_fields", srcStructWithFields, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
				DeclType: token.STRUCT,
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
		{"struct_with_fields_and_comments", srcStructWithFieldsAndComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
				DeclType: token.STRUCT,
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
		{"struct_block_decl", srcStructBlock, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
				DeclType: token.STRUCT,
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
			"filename:testfile.go+++pkg:main+++declType:struct+++name:picker": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:picker",
				DeclType: token.STRUCT,
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
		{"struct_with_hardwords", srcStructWithHardwords, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:struct+++name:httpClient": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:httpClient",
				DeclType: token.STRUCT,
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

			m := miner.NewDeclaration(lists.Dictionary)
			m.SetCurrentFile("testfile.go")
			ast.Walk(m, node)

			decls := m.Declarations()
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
		expected map[string]entity.Decl
	}{
		{"empty_interface", srcInterfaceWithoutMethods, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"selector": struct{}{},
				},
				Phrases: map[string]struct{}{}},
		}},
		{"interface_with_comments", srcInterfaceWithComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"comment":   struct{}{},
					"interface": struct{}{},
					"selector":  struct{}{},
				},
				Phrases: map[string]struct{}{
					"interface comment": struct{}{},
				}},
		}},
		{"interface_with_methods", srcInterfaceWithMethods, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"comment":   struct{}{},
					"interface": struct{}{},
					"pick":      struct{}{},
					"selector":  struct{}{},
				},
				Phrases: map[string]struct{}{
					"interface comment": struct{}{},
				}},
		}},
		{"interface_with_methods_and_comments", srcInterfaceWithMethodsAndComments, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"comment":   struct{}{},
					"function":  struct{}{},
					"interface": struct{}{},
					"pick":      struct{}{},
					"selector":  struct{}{},
				},
				Phrases: map[string]struct{}{
					"function comment":  struct{}{},
					"interface comment": struct{}{},
				}},
		}},
		{"interface_block_decl", srcInterfaceBlock, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
				DeclType: token.INTERFACE,
				Words: map[string]struct{}{
					"comment":   struct{}{},
					"function":  struct{}{},
					"global":    struct{}{},
					"interface": struct{}{},
					"pick":      struct{}{},
					"selector":  struct{}{},
				},
				Phrases: map[string]struct{}{
					"function comment":  struct{}{},
					"global comment":    struct{}{},
					"interface comment": struct{}{},
				}},
			"filename:testfile.go+++pkg:main+++declType:interface+++name:picker": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:picker",
				DeclType: token.INTERFACE,
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
		{"interface_with_hardwords", srcInterfaceWithHardwords, map[string]entity.Decl{
			"filename:testfile.go+++pkg:main+++declType:interface+++name:httpClient": entity.Decl{
				ID:       "filename:testfile.go+++pkg:main+++declType:interface+++name:httpClient",
				DeclType: token.INTERFACE,
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

			m := miner.NewDeclaration(lists.Dictionary)
			m.SetCurrentFile("testfile.go")
			ast.Walk(m, node)

			decls := m.Declarations()
			assert.Equal(t, len(fixture.expected), len(decls))
			assert.Equal(t, fixture.expected, decls)
		})
	}
}
