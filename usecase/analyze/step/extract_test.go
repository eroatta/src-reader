package step_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/usecase/analyze/step"
	"github.com/stretchr/testify/assert"
)

func TestExtract_OnNoFiles_ShouldReturnZeroIdentifiers(t *testing.T) {
	identc := step.Extract([]code.File{}, newExtractor)

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
	identc := step.Extract([]code.File{fileWithoutAST}, newExtractor)

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

	identc := step.Extract([]code.File{file}, newExtractor)

	identifiers := make(map[string]code.Identifier)
	for ident := range identc {
		identifiers[ident.Name] = ident
	}

	assert.Equal(t, 1, len(identifiers))
	assert.Equal(t, "main", identifiers["main"].Name)
}

func newExtractor(filename string) entity.Extractor {
	return &testExtractor{
		idents: make([]code.Identifier, 0),
	}
}

type testExtractor struct {
	idents []code.Identifier
}

func (t *testExtractor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch elem := node.(type) {
	case *ast.File:
		t.idents = append(t.idents, code.Identifier{
			Name:     elem.Name.String(),
			Position: elem.Pos(),
		})
	}

	return t
}

func (t *testExtractor) Identifiers() []code.Identifier {
	return t.idents
}
