package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewScopes_ShouldReturnScopesMiner(t *testing.T) {
	miner := miner.NewScope("testfile")

	assert.NotNil(t, miner)
}

func TestGetName_OnScope_ShouldReturnScope(t *testing.T) {
	miner := miner.NewScope("testfile")

	assert.Equal(t, "scope", miner.Name())
}

func TestScopedDeclarations_OnScope_ShouldReturnScopes(t *testing.T) {
	miner := miner.NewScope("testfile")

	assert.Equal(t, 0, len(miner.ScopedDeclarations()))
}

func TestVisit_OnScopeWithPlainFuncDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main

		func main() {

		}
	`

	expected := map[string]miner.ScopedDecl{
		"main++func::main": miner.ScopedDecl{
			ID:              "main++func::main",
			DeclType:        token.FUNC,
			Name:            "main",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        make([]string, 0),
			Comments:        make([]string, 0),
			PackageComments: make([]string, 0),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithFuncDeclWithComments_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		// package comment line 1
		// package comment line 2
		package main

		// function comment
		func main() {

		}
	`

	expected := map[string]miner.ScopedDecl{
		"main++func::main": miner.ScopedDecl{
			ID:            "main++func::main",
			DeclType:      token.FUNC,
			Name:          "main",
			VariableDecls: make([]string, 0),
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"function comment",
			},
			PackageComments: []string{
				"package comment line 1",
				"package comment line 2",
			},
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithMultipleFuncDeclWithComments_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		// package comment line 1
		// package comment line 2
		package main

		// function comment
		func main() {

		}

		// another comment
		func another() {

		}
	`

	expected := map[string]miner.ScopedDecl{
		"main++func::main": miner.ScopedDecl{
			ID:            "main++func::main",
			DeclType:      token.FUNC,
			Name:          "main",
			VariableDecls: make([]string, 0),
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"function comment",
			},
			PackageComments: []string{
				"package comment line 1",
				"package comment line 2",
			},
		},
		"main++func::another": miner.ScopedDecl{
			ID:            "main++func::another",
			DeclType:      token.FUNC,
			Name:          "another",
			VariableDecls: make([]string, 0),
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"another comment",
			},
			PackageComments: []string{
				"package comment line 1",
				"package comment line 2",
			},
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithMultipleFuncDeclWithFullBody_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		// package comment line 1
		// package comment line 2
		package main

		// function comment
		func main() {

		}

		// another comment
		func another() {

		}
	`

	expected := map[string]miner.ScopedDecl{
		"main++func::main": miner.ScopedDecl{
			ID:            "main++func::main",
			DeclType:      token.FUNC,
			Name:          "main",
			VariableDecls: make([]string, 0),
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"function comment",
			},
			PackageComments: []string{
				"package comment line 1",
				"package comment line 2",
			},
		},
		"main++func::another": miner.ScopedDecl{
			ID:            "main++func::another",
			DeclType:      token.FUNC,
			Name:          "another",
			VariableDecls: make([]string, 0),
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"another comment",
			},
			PackageComments: []string{
				"package comment line 1",
				"package comment line 2",
			},
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithPlainVarDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main

		var Common string
	`

	expected := map[string]miner.ScopedDecl{
		"main++var::Common": miner.ScopedDecl{
			ID:              "main++var::Common",
			DeclType:        token.VAR,
			Name:            "Common",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        make([]string, 0),
			Comments:        make([]string, 0),
			PackageComments: make([]string, 0),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithFullyCommentedVarDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		// package comment
		package main

		// outer comment
		var Common string // inner comment
	`

	expected := map[string]miner.ScopedDecl{
		"main++var::Common": miner.ScopedDecl{
			ID:            "main++var::Common",
			DeclType:      token.VAR,
			Name:          "Common",
			VariableDecls: make([]string, 0),
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"outer comment",
				"inner comment",
			},
			PackageComments: []string{"package comment"},
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithVarBlockDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main
		
		// outer comment
		var (
			common string
			regular string = "valid"
			nrzXXZ int = 32
		)
	`

	expected := map[string]miner.ScopedDecl{
		"main++var::common": miner.ScopedDecl{
			ID:              "main++var::common",
			DeclType:        token.VAR,
			Name:            "common",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        make([]string, 0),
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"main++var::regular": miner.ScopedDecl{
			ID:              "main++var::regular",
			DeclType:        token.VAR,
			Name:            "regular",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"valid"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"main++var::nrzXXZ": miner.ScopedDecl{
			ID:              "main++var::nrzXXZ",
			DeclType:        token.VAR,
			Name:            "nrzXXZ",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        make([]string, 0),
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithPlainConstDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main

		const Common string = "valid"
	`

	expected := map[string]miner.ScopedDecl{
		"main++const::Common": miner.ScopedDecl{
			ID:              "main++const::Common",
			DeclType:        token.CONST,
			Name:            "Common",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"valid"},
			Comments:        make([]string, 0),
			PackageComments: make([]string, 0),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithFullyCommentedConstDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		// package comment
		package main

		// outer comment
		const Common string = "valid" // inner comment
	`

	expected := map[string]miner.ScopedDecl{
		"main++const::Common": miner.ScopedDecl{
			ID:            "main++const::Common",
			DeclType:      token.CONST,
			Name:          "Common",
			VariableDecls: make([]string, 0),
			Statements:    make([]string, 0),
			BodyText:      []string{"valid"},
			Comments: []string{
				"outer comment",
				"inner comment",
			},
			PackageComments: []string{"package comment"},
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithConstBlockDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main
		
		// outer comment
		const (
			common string = "common"
			regular, notRegular string = "valid", "invalid"
			nrzXXZ int = 32
		)
	`

	expected := map[string]miner.ScopedDecl{
		"main++const::common": miner.ScopedDecl{
			ID:              "main++const::common",
			DeclType:        token.CONST,
			Name:            "common",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"common"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"main++const::regular": miner.ScopedDecl{
			ID:              "main++const::regular",
			DeclType:        token.CONST,
			Name:            "regular",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"valid"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"main++const::notRegular": miner.ScopedDecl{
			ID:              "main++const::notRegular",
			DeclType:        token.CONST,
			Name:            "notRegular",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"invalid"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"main++const::nrzXXZ": miner.ScopedDecl{
			ID:              "main++const::nrzXXZ",
			DeclType:        token.CONST,
			Name:            "nrzXXZ",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        make([]string, 0),
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithPlainStructDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main

		type selector struct{}
	`

	expected := map[string]miner.ScopedDecl{
		"main++struct::selector": miner.ScopedDecl{
			ID:              "main++struct::selector",
			DeclType:        token.STRUCT,
			Name:            "selector",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        make([]string, 0),
			Comments:        make([]string, 0),
			PackageComments: make([]string, 0),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithFullyCommentedStructDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		// package comment
		package main

		// type comment
		type selector struct{
			// field comment
			pick string
		}
	`

	expected := map[string]miner.ScopedDecl{
		"main++struct::selector": miner.ScopedDecl{
			ID:            "main++struct::selector",
			DeclType:      token.STRUCT,
			Name:          "selector",
			VariableDecls: []string{"pick string"},
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"field comment",
				"type comment",
			},
			PackageComments: []string{"package comment"},
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithStructBlockDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		// package comment
		package main

		// global comment
		type (
			// local comment
			selector struct {
				// field comment
				pick string
			}

			// inner comment
			httpClient struct {
				protocolPicker string
				url string
			}
		)
	`

	expected := map[string]miner.ScopedDecl{
		"main++struct::selector": miner.ScopedDecl{
			ID:            "main++struct::selector",
			DeclType:      token.STRUCT,
			Name:          "selector",
			VariableDecls: []string{"pick string"},
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"field comment",
				"local comment",
			},
			PackageComments: []string{"package comment"},
		},
		"main++struct::httpClient": miner.ScopedDecl{
			ID:       "main++struct::httpClient",
			DeclType: token.STRUCT,
			Name:     "http client",
			VariableDecls: []string{
				"protocol picker string",
				"url string",
			},
			Statements: make([]string, 0),
			BodyText:   make([]string, 0),
			Comments: []string{
				"global comment",
				"inner comment",
			},
			PackageComments: []string{"package comment"},
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	m := miner.NewScope("testfile")
	ast.Walk(m, node)

	scopedDecls := m.ScopedDeclarations()
	assert.Equal(t, expected, scopedDecls)
}

/*
func TestVisit_OnScopeWithInterfaceDecl_ShouldReturnDeclarationScopesForInterfaces(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}*/
