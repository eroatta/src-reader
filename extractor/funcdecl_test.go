package extractor_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/extractor"
	"github.com/stretchr/testify/assert"
)

func TestNodeType_OnFuncDecl_ShouldReturnFuncDecl(t *testing.T) {
	extr := extractor.FuncDecl{}

	assert.Equal(t, "*ast.FuncDecl", extr.NodeType().String())
}

func TestExtract_OnFuncDeclWithNoFuncDeclNode_ShouldReturnEmptyArray(t *testing.T) {
	testFileset := token.NewFileSet()
	ast, _ := parser.ParseFile(testFileset, "main.go", `package main`, parser.AllErrors)
	extr := extractor.FuncDecl{}

	identifiers := extr.Extract("test.go", ast)

	assert.Empty(t, identifiers)
}

func TestExtract_OnFuncDeclWithFuncDeclNode_ShouldReturnOneElement(t *testing.T) {
	src := `
		package main

		func main() {
			// empty func
		}
	`
	testFileset := token.NewFileSet()
	tree, _ := parser.ParseFile(testFileset, "test.go", src, parser.AllErrors)

	identifiers := extractor.FuncDecl{}.Extract("test.go", tree.Decls[0])

	assert.Equal(t, 1, len(identifiers))

	retrievedIdent := identifiers[0]
	assert.Equal(t, "main", retrievedIdent.Name)
	assert.Equal(t, "test.go", retrievedIdent.File)
	assert.Equal(t, "FuncDecl", retrievedIdent.Type)

	funcDecl, _ := tree.Decls[0].(*ast.FuncDecl)
	assert.Equal(t, funcDecl.Name.Pos(), retrievedIdent.Position)
}
