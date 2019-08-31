package step_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/step"
	"github.com/stretchr/testify/assert"
)

func TestExtract_OnNoFiles_ShouldReturnZeroIdentifiers(t *testing.T) {
	identc := step.Extract([]code.File{}, extractor{})

	var identifiers int
	for range identc {
		identifiers++
	}

	assert.Equal(t, 0, identifiers)
}

func TestExtract_OnFileWithoutAST_ShouldReturnZeroIdentifiers(t *testing.T) {
	fileWithoutAST := code.File{
		Name: "main.go",
		AST:  nil,
	}
	identc := step.Extract([]code.File{fileWithoutAST}, extractor{})

	var identifiers int
	for range identc {
		identifiers++
	}

	assert.Equal(t, 0, identifiers)
}

func TestExtract_OnFileWithAST_ShouldReturnFoundIdentifiers(t *testing.T) {
	/* Created AST:
	    0  *ast.File {
	    1  .  Doc: nil
	    2  .  Package: 1:1
	    3  .  Name: *ast.Ident {
	    4  .  .  NamePos: 1:9
	    5  .  .  Name: "main"
	    6  .  .  Obj: nil
	    7  .  }
	    8  .  Decls: nil
	    9  .  Scope: *ast.Scope {
	   10  .  .  Outer: nil
	   11  .  .  Objects: map[string]*ast.Object (len = 0) {}
	   12  .  }
	   13  .  Imports: nil
	   14  .  Unresolved: nil
	   15  .  Comments: nil
	   16  }
	*/
	testFileset := token.NewFileSet()

	ast, _ := parser.ParseFile(testFileset, "main.go", `package main`, parser.AllErrors)
	file := code.File{
		Name:    "main.go",
		AST:     ast,
		FileSet: testFileset,
	}

	identc := step.Extract([]code.File{file}, extractor{})

	identifiers := make(map[string]code.Identifier)
	for ident := range identc {
		identifiers[ident.Name] = ident
	}

	assert.Equal(t, 1, len(identifiers))
	assert.Equal(t, "main", identifiers["main"].Name)
}

type extractor struct {
	node *ast.Ident
}

func (e extractor) NodeType() reflect.Type {
	return reflect.TypeOf(e.node)
}

func (e extractor) Extract(filename string, node ast.Node) []code.Identifier {
	i, ok := node.(*ast.Ident)
	if ok {
		ident := code.Identifier{
			Name:     i.Name,
			Position: i.Pos(),
		}

		return []code.Identifier{ident}
	}

	return []code.Identifier{}
}
