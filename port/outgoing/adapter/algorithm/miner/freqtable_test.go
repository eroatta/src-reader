package miner_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/token/samurai"

	"github.com/eroatta/src-reader/port/outgoing/adapter/algorithm/miner"
	"github.com/stretchr/testify/assert"
)

func TestNewGlobalFreqTableFactory_ShouldReturnGlobalFreqTableMinerFactory(t *testing.T) {
	factory := miner.NewGlobalFreqTableFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnGlobalFreqTableFactory_ShouldReturnMiner(t *testing.T) {
	factory := miner.NewGlobalFreqTableFactory()
	miner, err := factory.Make()

	assert.Equal(t, "global-frequency-table", miner.Name())
	assert.NoError(t, err)
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

	gft := miner.Results().(*samurai.FrequencyTable)
	assert.Equal(t, ft, gft)
}
