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

func TestVisit_OnSamuraiWithVarDeclBlock_ShouldSplitAllTheNames(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main
		var (
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

func TestVisit_OnSamuraiWithFuncDeclNode_ShouldSplitTheName(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main

		import "fmt"

		func main() {
			fmt.Println("Hello Samurai Extractor!")
		}
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["main"])
}

func TestVisit_OnSamuraiWithFuncDeclNodeWithArgs_ShouldSplitAllTheNames(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main

		import (
			"fmt"
			"gin"
		)

		func main() {
			engine := gin.New()
			engine.Delims("first_argument", "second_argument")
		}

		func (engine *Engine) Delims(left, right string) *Engine {
			engine.delims = render.Delims{Left: left, Right: right}
			return engine
		}
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["main"])
	assert.Equal(t, 1, extractor.words["delims"])
	assert.Equal(t, 1, extractor.words["left"])
	assert.Equal(t, 1, extractor.words["right"])
}

func TestVisit_OnSamuraiWithFuncDeclNodeWithNamedResults_ShouldSplitAllTheNames(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main

		import (
			"fmt"
		)
		
		func main() {
			text, number := results()
			fmt.Println(fmt.Sprintf("%s, %d", text, number))
		}
		
		func results() (text string, number int) {
			return "text", 10
		}
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["main"])
	assert.Equal(t, 1, extractor.words["results"])
	assert.Equal(t, 1, extractor.words["text"])
	assert.Equal(t, 1, extractor.words["number"])
}

func TestVisit_OnSamuraiWithTypeDeclNode_ShouldSplitTheName(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main

		type myInt int
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["my"])
	assert.Equal(t, 1, extractor.words["int"])
}

func TestVisit_OnSamuraiWithStructTypeDeclNodeWithNoFields_ShouldSplitTheName(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main

		type myStruct struct {}
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["my"])
	assert.Equal(t, 1, extractor.words["struct"])
}

func TestVisit_OnSamuraiWithStructTypeDeclNodeWithTwoFields_ShouldSplitAllTheNames(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main

		type myStruct struct {
			first string,
			second string
		}
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["my"])
	assert.Equal(t, 1, extractor.words["struct"])
	assert.Equal(t, 1, extractor.words["first"])
	assert.Equal(t, 1, extractor.words["second"])
}

func TestVisit_OnSamuraiWithInterfaceTypeDeclNode_ShouldSplitTheName(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main

		type myInterface interface
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 1, extractor.words["my"])
	assert.Equal(t, 1, extractor.words["interface"])
}

func TestVisit_OnSamuraiWithInterfaceTypeDeclNodeWithoutMethods_ShouldSplitTheName(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		package main

		type myInterface interface {
			myMethod(arg string) (out string)
		}
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.AllErrors)
	//ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 2, extractor.words["my"])
	assert.Equal(t, 1, extractor.words["interface"])
	assert.Equal(t, 1, extractor.words["method"])
	assert.Equal(t, 1, extractor.words["out"])
}

func TestVisit_OnSamuraiWithCommentsOnFile_ShouldSplitTheComments(t *testing.T) {
	samurai := NewSamuraiExtractor()

	fs := token.NewFileSet()
	var src = `
		// package comment
		package main
	
		// regular comment
		type MyInterface interface {
			MyMethod(out string) // inline comment
		}
	
		/* 
		block comment
		*/
		type abc int
	`

	node, _ := parser.ParseFile(fs, "", []byte(src), parser.ParseComments)
	ast.Print(fs, node)
	ast.Walk(samurai, node)

	assert.NotNil(t, samurai)

	extractor := samurai.(SamuraiExtractor)
	assert.NotEmpty(t, extractor.words)
	assert.Equal(t, 2, extractor.words["my"])
	assert.Equal(t, 1, extractor.words["interface"])
	assert.Equal(t, 1, extractor.words["method"])
	assert.Equal(t, 1, extractor.words["out"])
	assert.Equal(t, 1, extractor.words["abc"])
	assert.Equal(t, 1, extractor.words["package"])
	assert.Equal(t, 1, extractor.words["regular"])
	assert.Equal(t, 1, extractor.words["inline"])
	assert.Equal(t, 1, extractor.words["block"])
	assert.Equal(t, 4, extractor.words["comment"])
}
