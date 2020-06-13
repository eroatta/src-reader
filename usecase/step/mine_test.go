package step_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/usecase/step"
	"github.com/stretchr/testify/assert"
)

func TestMine_OnNoFiles_ShouldReturnMinersWithoutResults(t *testing.T) {
	processed := step.Mine([]entity.File{}, &miner{typ: "empty"})

	assert.Equal(t, 1, len(processed))

	emptyMiner, ok := processed["empty"].(*miner)
	assert.True(t, ok)
	assert.NotNil(t, emptyMiner)
	assert.Equal(t, 0, emptyMiner.visits)
}

func TestMine_OnEmptyMiners_ShouldReturnNoResults(t *testing.T) {
	processed := step.Mine([]entity.File{}, []entity.Miner{}...)

	assert.Equal(t, 0, len(processed))
}

func TestMine_OnFileWithNilAST_ShouldReturnMinersWithoutResults(t *testing.T) {
	processed := step.Mine([]entity.File{{Name: "main.go"}}, &miner{typ: "empty"})

	assert.Equal(t, 1, len(processed))

	emptyMiner, ok := processed["empty"].(*miner)
	assert.True(t, ok)
	assert.NotNil(t, emptyMiner)
	assert.Equal(t, 0, emptyMiner.visits)
}

func TestMine_OnTwoMiners_ShouldReturnResultsBothMiners(t *testing.T) {
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

	ast1, _ := parser.ParseFile(testFileset, "main.go", `package main`, parser.AllErrors)
	file1 := entity.File{
		Name:    "main.go",
		AST:     ast1,
		FileSet: testFileset,
	}

	ast2, _ := parser.ParseFile(testFileset, "test.go", `package test`, parser.AllErrors)
	file2 := entity.File{
		Name:    "test.go",
		AST:     ast2,
		FileSet: testFileset,
	}

	first := &miner{typ: "first"}
	second := &miner{typ: "second"}

	processed := step.Mine([]entity.File{file1, file2}, first, second)

	assert.Equal(t, 2, len(processed))

	firstMiner, ok := processed["first"].(*miner)
	assert.True(t, ok)
	assert.NotNil(t, firstMiner)
	assert.Equal(t, 8, firstMiner.visits)

	secondMiner, ok := processed["second"].(*miner)
	assert.True(t, ok)
	assert.NotNil(t, secondMiner)
	assert.Equal(t, 8, secondMiner.visits)
}

type miner struct {
	typ    string
	visits int
}

func (m *miner) Name() string {
	return m.typ
}

func (m *miner) SetCurrentFile(f string) {}

func (m *miner) Visit(n ast.Node) ast.Visitor {
	m.visits++
	return m
}

func (m *miner) Results() interface{} {
	return m.visits
}
