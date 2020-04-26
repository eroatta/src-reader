package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/adapter/algorithm/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewCommentsFactory_ShouldReturnCommentsMinerFactory(t *testing.T) {
	factory := miner.NewCommentsFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnCommentsFactory_ShouldReturnMiner(t *testing.T) {
	factory := miner.NewCommentsFactory()
	miner, err := factory.Make()

	assert.Equal(t, "comments", miner.Name())
	assert.NoError(t, err)
}

func TestVisit_OnComments_ShouldReturnCleanComments(t *testing.T) {
	src := `
		// package comment
		// package comment 2
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

	expected := []string{
		"package comment",
		"package comment 2",
		"global comment",
		"interface comment",
		"function comment",
		"inner comment",
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	f := miner.NewCommentsFactory()
	c, _ := f.Make()
	ast.Walk(c, node)

	collected := c.Results().([]string)
	assert.ElementsMatch(t, expected, collected)
}
