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
			// inner comment
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

/*func TestVisit_OnScopeWithVarDecl_ShouldReturnDeclarationScopesForVars(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestVisit_OnScopeWithConstDecl_ShouldReturnDeclarationScopesForConsts(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestVisit_OnScopeWithStructDecl_ShouldReturnDeclarationScopesForStructs(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestVisit_OnScopeWithInterfaceDecl_ShouldReturnDeclarationScopesForInterfaces(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}*/
