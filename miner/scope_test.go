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

func TestVisit_OnScopeWithFuncDecl_ShouldReturnDeclarationScopesForFuncs(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected map[string]miner.ScopedDecl
	}{}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, _ := parser.ParseFile(fs, "testfile", []byte(fixture.src), parser.ParseComments)

			m := miner.NewScope("testfile")
			ast.Walk(m, node)

			scopedDecls := m.ScopedDeclarations()
			assert.Equal(t, fixture.expected, scopedDecls)
		})
	}

	assert.FailNow(t, "not yet implemented")
}

func TestVisit_OnScopeWithVarDecl_ShouldReturnDeclarationScopesForVars(t *testing.T) {
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
}
