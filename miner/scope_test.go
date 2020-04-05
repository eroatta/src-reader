package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewScopes_ShouldReturnScopesMiner(t *testing.T) {
	miner := miner.NewScope("testfile")

	assert.NotNil(t, miner)
}

func TestGetName_OnScope_ShouldReturnScope(t *testing.T) {
	miner := miner.NewScope("testfile")

	assert.Equal(t, entity.MinerScopedDeclarations, miner.Type())
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

	expected := map[string]entity.ScopedDecl{
		"main++func::main": entity.ScopedDecl{
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

	expected := map[string]entity.ScopedDecl{
		"main++func::main": entity.ScopedDecl{
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

	expected := map[string]entity.ScopedDecl{
		"main++func::main": entity.ScopedDecl{
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
		"main++func::another": entity.ScopedDecl{
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
			// comment inside another
		}

		func (m Miner) apply(strategy string, nodes int64) (results []string) {
			// apply strategy comment
			results = append(results, "result 1")
			return results
		}
	`

	expected := map[string]entity.ScopedDecl{
		"main++func::main": entity.ScopedDecl{
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
		"main++func::another": entity.ScopedDecl{
			ID:            "main++func::another",
			DeclType:      token.FUNC,
			Name:          "another",
			VariableDecls: make([]string, 0),
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"another comment",
				"comment inside another",
			},
			PackageComments: []string{
				"package comment line 1",
				"package comment line 2",
			},
		},
		"main++func::apply": entity.ScopedDecl{
			ID:       "main++func::apply",
			DeclType: token.FUNC,
			Name:     "apply",
			VariableDecls: []string{
				"strategy string",
				"nodes int64",
				"results []string",
			},
			Statements: make([]string, 0),
			BodyText:   make([]string, 0),
			Comments: []string{
				"apply strategy comment",
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

	expected := map[string]entity.ScopedDecl{
		"main++var::Common": entity.ScopedDecl{
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

	expected := map[string]entity.ScopedDecl{
		"main++var::Common": entity.ScopedDecl{
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

	expected := map[string]entity.ScopedDecl{
		"main++var::common": entity.ScopedDecl{
			ID:              "main++var::common",
			DeclType:        token.VAR,
			Name:            "common",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        make([]string, 0),
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"main++var::regular": entity.ScopedDecl{
			ID:              "main++var::regular",
			DeclType:        token.VAR,
			Name:            "regular",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"valid"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"main++var::nrzXXZ": entity.ScopedDecl{
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

	expected := map[string]entity.ScopedDecl{
		"main++const::Common": entity.ScopedDecl{
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

	expected := map[string]entity.ScopedDecl{
		"main++const::Common": entity.ScopedDecl{
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

	expected := map[string]entity.ScopedDecl{
		"main++const::common": entity.ScopedDecl{
			ID:              "main++const::common",
			DeclType:        token.CONST,
			Name:            "common",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"common"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"main++const::regular": entity.ScopedDecl{
			ID:              "main++const::regular",
			DeclType:        token.CONST,
			Name:            "regular",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"valid"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"main++const::notRegular": entity.ScopedDecl{
			ID:              "main++const::notRegular",
			DeclType:        token.CONST,
			Name:            "notRegular",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"invalid"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"main++const::nrzXXZ": entity.ScopedDecl{
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

	expected := map[string]entity.ScopedDecl{
		"main++struct::selector": entity.ScopedDecl{
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

	expected := map[string]entity.ScopedDecl{
		"main++struct::selector": entity.ScopedDecl{
			ID:            "main++struct::selector",
			DeclType:      token.STRUCT,
			Name:          "selector",
			VariableDecls: []string{"pick string"},
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"type comment",
				"field comment",
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

	expected := map[string]entity.ScopedDecl{
		"main++struct::selector": entity.ScopedDecl{
			ID:            "main++struct::selector",
			DeclType:      token.STRUCT,
			Name:          "selector",
			VariableDecls: []string{"pick string"},
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"global comment",
				"local comment",
				"field comment",
			},
			PackageComments: []string{"package comment"},
		},
		"main++struct::httpClient": entity.ScopedDecl{
			ID:       "main++struct::httpClient",
			DeclType: token.STRUCT,
			Name:     "httpClient",
			VariableDecls: []string{
				"protocolPicker string",
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

func TestVisit_OnScopeWithPlainInterfaceDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main

		type selector interface{}
	`

	expected := map[string]entity.ScopedDecl{
		"main++interface::selector": entity.ScopedDecl{
			ID:              "main++interface::selector",
			DeclType:        token.INTERFACE,
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

func TestVisit_OnScopeWithFullyCommentedInterfaceDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		// package comment
		package main

		// type comment
		type selector interface{
			// function comment
			pick() string
		}
	`

	expected := map[string]entity.ScopedDecl{
		"main++interface::selector": entity.ScopedDecl{
			ID:            "main++interface::selector",
			DeclType:      token.INTERFACE,
			Name:          "selector",
			VariableDecls: make([]string, 0),
			Statements:    []string{"pick"},
			BodyText:      make([]string, 0),
			Comments: []string{
				"type comment",
				"function comment",
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

func TestVisit_OnScopeWithInterfaceBlockDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		// package comment
		package main

		// global comment
		type (
			// interface comment
			selector interface {
				// function comment
				pick() string
			}

			// inner comment
			httpClient interface {
				protocolPicker() string
				url() string
			}
		)
	`

	expected := map[string]entity.ScopedDecl{
		"main++interface::selector": entity.ScopedDecl{
			ID:            "main++interface::selector",
			DeclType:      token.INTERFACE,
			Name:          "selector",
			VariableDecls: make([]string, 0),
			Statements:    []string{"pick"},
			BodyText:      make([]string, 0),
			Comments: []string{
				"global comment",
				"interface comment",
				"function comment",
			},
			PackageComments: []string{"package comment"},
		},
		"main++interface::httpClient": entity.ScopedDecl{
			ID:            "main++interface::httpClient",
			DeclType:      token.INTERFACE,
			Name:          "httpClient",
			VariableDecls: make([]string, 0),
			Statements: []string{
				"protocolPicker",
				"url",
			},
			BodyText: make([]string, 0),
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
