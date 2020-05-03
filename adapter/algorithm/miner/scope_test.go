package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/adapter/algorithm/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewScopesFactory_ShouldReturnScopesMinerFactory(t *testing.T) {
	factory := miner.NewScopesFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnScopesFactory_ShouldReturnMiner(t *testing.T) {
	factory := miner.NewScopesFactory()
	miner, err := factory.Make()

	assert.Equal(t, "scoped-declarations", miner.Name())
	assert.NoError(t, err)
}

func TestSetCurrentFile_OnScope_ShoudldSetNewFilename(t *testing.T) {
	miner := miner.NewScope()
	miner.SetCurrentFile("new_file.go")

	assert.Equal(t, "new_file.go", miner.Filename)
}

func TestScopedDeclarations_OnScope_ShouldReturnScopes(t *testing.T) {
	mnr := miner.NewScope()

	results := mnr.Results().(map[string]miner.ScopedDecl)
	assert.Equal(t, 0, len(results))
}

func TestVisit_OnScopeWithPlainFuncDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main

		func main() {

		}
	`

	expected := map[string]miner.ScopedDecl{
		"filename:testfile.go+++pkg:main+++declType:func+++name:main": {
			ID:              "filename:testfile.go+++pkg:main+++declType:func+++name:main",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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
		"filename:testfile.go+++pkg:main+++declType:func+++name:main": {
			ID:            "filename:testfile.go+++pkg:main+++declType:func+++name:main",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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
		"filename:testfile.go+++pkg:main+++declType:func+++name:main": {
			ID:            "filename:testfile.go+++pkg:main+++declType:func+++name:main",
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
		"filename:testfile.go+++pkg:main+++declType:func+++name:another": {
			ID:            "filename:testfile.go+++pkg:main+++declType:func+++name:another",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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

	expected := map[string]miner.ScopedDecl{
		"filename:testfile.go+++pkg:main+++declType:func+++name:main": {
			ID:            "filename:testfile.go+++pkg:main+++declType:func+++name:main",
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
		"filename:testfile.go+++pkg:main+++declType:func+++name:another": {
			ID:            "filename:testfile.go+++pkg:main+++declType:func+++name:another",
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
		"filename:testfile.go+++pkg:main+++declType:func+++name:Miner.apply": {
			ID:       "filename:testfile.go+++pkg:main+++declType:func+++name:Miner.apply",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithPlainVarDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main

		var Common string
	`

	expected := map[string]miner.ScopedDecl{
		"filename:testfile.go+++pkg:main+++declType:var+++name:Common": {
			ID:              "filename:testfile.go+++pkg:main+++declType:var+++name:Common",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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
		"filename:testfile.go+++pkg:main+++declType:var+++name:Common": {
			ID:            "filename:testfile.go+++pkg:main+++declType:var+++name:Common",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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
		"filename:testfile.go+++pkg:main+++declType:var+++name:common": {
			ID:              "filename:testfile.go+++pkg:main+++declType:var+++name:common",
			DeclType:        token.VAR,
			Name:            "common",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        make([]string, 0),
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"filename:testfile.go+++pkg:main+++declType:var+++name:regular": {
			ID:              "filename:testfile.go+++pkg:main+++declType:var+++name:regular",
			DeclType:        token.VAR,
			Name:            "regular",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"valid"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"filename:testfile.go+++pkg:main+++declType:var+++name:nrzXXZ": {
			ID:              "filename:testfile.go+++pkg:main+++declType:var+++name:nrzXXZ",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithPlainConstDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main

		const Common string = "valid"
	`

	expected := map[string]miner.ScopedDecl{
		"filename:testfile.go+++pkg:main+++declType:const+++name:Common": {
			ID:              "filename:testfile.go+++pkg:main+++declType:const+++name:Common",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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
		"filename:testfile.go+++pkg:main+++declType:const+++name:Common": {
			ID:            "filename:testfile.go+++pkg:main+++declType:const+++name:Common",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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
		"filename:testfile.go+++pkg:main+++declType:const+++name:common": {
			ID:              "filename:testfile.go+++pkg:main+++declType:const+++name:common",
			DeclType:        token.CONST,
			Name:            "common",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"common"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"filename:testfile.go+++pkg:main+++declType:const+++name:regular": {
			ID:              "filename:testfile.go+++pkg:main+++declType:const+++name:regular",
			DeclType:        token.CONST,
			Name:            "regular",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"valid"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"filename:testfile.go+++pkg:main+++declType:const+++name:notRegular": {
			ID:              "filename:testfile.go+++pkg:main+++declType:const+++name:notRegular",
			DeclType:        token.CONST,
			Name:            "notRegular",
			VariableDecls:   make([]string, 0),
			Statements:      make([]string, 0),
			BodyText:        []string{"invalid"},
			Comments:        []string{"outer comment"},
			PackageComments: make([]string, 0),
		},
		"filename:testfile.go+++pkg:main+++declType:const+++name:nrzXXZ": {
			ID:              "filename:testfile.go+++pkg:main+++declType:const+++name:nrzXXZ",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithPlainStructDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main

		type selector struct{}
	`

	expected := map[string]miner.ScopedDecl{
		"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": {
			ID:              "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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
		"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": {
			ID:            "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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

		// func on struct comment
		func (s selector) print() {}
	`

	expected := map[string]miner.ScopedDecl{
		"filename:testfile.go+++pkg:main+++declType:struct+++name:selector": {
			ID:            "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
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
		"filename:testfile.go+++pkg:main+++declType:struct+++name:httpClient": {
			ID:       "filename:testfile.go+++pkg:main+++declType:struct+++name:httpClient",
			DeclType: token.STRUCT,
			Name:     "httpClient",
			VariableDecls: []string{
				"protocolpicker string",
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
		"filename:testfile.go+++pkg:main+++declType:func+++name:selector.print": {
			ID:            "filename:testfile.go+++pkg:main+++declType:func+++name:selector.print",
			DeclType:      token.FUNC,
			Name:          "print",
			VariableDecls: []string{},
			Statements:    make([]string, 0),
			BodyText:      make([]string, 0),
			Comments: []string{
				"func on struct comment",
			},
			PackageComments: []string{"package comment"},
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
	assert.Equal(t, expected, scopedDecls)
}

func TestVisit_OnScopeWithPlainInterfaceDecl_ShouldReturnScopedDeclaration(t *testing.T) {
	src := `
		package main

		type selector interface{}
	`

	expected := map[string]miner.ScopedDecl{
		"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": {
			ID:              "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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

	expected := map[string]miner.ScopedDecl{
		"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": {
			ID:            "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
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

	expected := map[string]miner.ScopedDecl{
		"filename:testfile.go+++pkg:main+++declType:interface+++name:selector": {
			ID:            "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
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
		"filename:testfile.go+++pkg:main+++declType:interface+++name:httpClient": {
			ID:            "filename:testfile.go+++pkg:main+++declType:interface+++name:httpClient",
			DeclType:      token.INTERFACE,
			Name:          "httpClient",
			VariableDecls: make([]string, 0),
			Statements: []string{
				"protocolpicker",
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
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	m := miner.NewScope()
	m.SetCurrentFile("testfile.go")
	ast.Walk(m, node)

	scopedDecls := m.Results().(map[string]miner.ScopedDecl)
	assert.Equal(t, expected, scopedDecls)
}
