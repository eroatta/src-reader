package extractors

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: use a real test... not this piece of shit
func TestProcess_OnSamuraiExtractor_ShouldNotFail(t *testing.T) {
	samurai := NewSamuraiExtractor()

	Process(samurai, nil)

	assert.True(t, true)
}

func TestNewSamurai_ShouldReturnNewExtractor(t *testing.T) {
	extractor := NewSamuraiExtractor()

	assert.NotNil(t, extractor)
	assert.IsType(t, SamuraiExtractor{}, extractor)
}

func TestGetName_OnSamurai_ShouldReturnSamurai(t *testing.T) {
	extractor := NewSamuraiExtractor()

	assert.Equal(t, "samurai", extractor.Name())
}

func TestVisit_OnSamuraiWithNilNode_ShouldReturnNil(t *testing.T) {
	extractor := NewSamuraiExtractor()

	assert.Nil(t, extractor.Visit(nil))
}

func TestVisit_OnSamuraiWithNotNilNode_ShouldReturnVisitor(t *testing.T) {
	extractor := NewSamuraiExtractor()

	node, _ := parser.ParseExpr("a + b")
	got := extractor.Visit(node)

	assert.NotNil(t, got)
}

func TestVisit_OnSamuraiWithVarDeclNode_ShouldSplitTheName(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main
		var testIden string
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["test"])
	assert.Equal(t, 1, extractor.words["iden"])
}

func TestVisit_OnSamuraiWithVarDeclAndMultipleNames_ShouldSplitAllTheNames(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main
		var testIden, testBar string
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 2, extractor.words["test"])
	assert.Equal(t, 1, extractor.words["iden"])
	assert.Equal(t, 1, extractor.words["bar"])
}

func TestVisit_OnSamuraiWithConstDeclNode_ShouldSplitTheName(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main
		const testIden = "boo"
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["test"])
	assert.Equal(t, 1, extractor.words["iden"])
	assert.Equal(t, 1, extractor.words["boo"])
}

func TestVisit_OnSamuraiWithConstDeclAndMultipleNamesAndValues_ShouldSplitAllTheNames(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main
		const testIden, foo = "boo", "bar"
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["test"])
	assert.Equal(t, 1, extractor.words["iden"])
	assert.Equal(t, 1, extractor.words["boo"])
	assert.Equal(t, 1, extractor.words["foo"])
	assert.Equal(t, 1, extractor.words["bar"])
}

func TestVisit_OnSamuraiWithConstDeclBlock_ShouldSplitAllTheNames(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main
		const (
			testIden = "boo"
			foo = "bar"
		)
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["test"])
	assert.Equal(t, 1, extractor.words["iden"])
	assert.Equal(t, 1, extractor.words["boo"])
	assert.Equal(t, 1, extractor.words["foo"])
	assert.Equal(t, 1, extractor.words["bar"])
}

func TestVisit_OnSamuraiWithTypeDeclNode_ShouldSplitTheName(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestVisit_OnSamuraiWithTypeDeclNodeAndMultipleNames_ShouldSplitAllTheNames(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestVisit_OnSamuraiWithFuncDeclNode_ShouldSplitTheName(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

// TODO: add tests for interfaces, imports and comments
