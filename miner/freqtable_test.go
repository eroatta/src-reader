package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/token/samurai"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewGlobalFreqTable_ShouldReturnNewGlobalFreqTableMiner(t *testing.T) {
	miner := miner.NewGlobalFreqTable(nil)

	assert.NotNil(t, miner)
}

func TestType_OnGlobalFreqTable_ShouldReturnMinerType(t *testing.T) {
	miner := miner.NewGlobalFreqTable(nil)

	assert.Equal(t, entity.MinerGlobalFrequencyTable, miner.Type())
}

func TestVisit_OnGlobalFreqTable_ShouldReturnCleanComments(t *testing.T) {
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

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	ft := samurai.NewFrequencyTable()
	miner := miner.NewGlobalFreqTable(ft)
	ast.Walk(miner, node)

	gft := miner.Table()
	assert.Equal(t, ft, gft)
}
