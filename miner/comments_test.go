package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewComments_ShouldReturnNewCommentsMiner(t *testing.T) {
	miner := miner.NewComments()

	assert.NotNil(t, miner)
}

func TestType_OnComments_ShouldReturnMinerType(t *testing.T) {
	miner := miner.NewComments()

	assert.Equal(t, entity.MinerComments, miner.Type())
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

	c := miner.NewComments()
	ast.Walk(c, node)

	collected := c.Collected()
	assert.ElementsMatch(t, expected, collected)
}
