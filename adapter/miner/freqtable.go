package miner

import (
	"go/ast"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/samurai"
)

type GlobalFreqTable struct {
	freqTable *samurai.FrequencyTable
}

func NewGlobalFreqTable(ft *samurai.FrequencyTable) GlobalFreqTable {
	return GlobalFreqTable{
		freqTable: ft,
	}
}

func (g GlobalFreqTable) Type() entity.MinerType {
	return entity.MinerGlobalFrequencyTable
}

func (g GlobalFreqTable) SetCurrentFile(filename string) {
	// do nothing
}

func (g GlobalFreqTable) Visit(node ast.Node) ast.Visitor {
	return nil
}

func (g GlobalFreqTable) Table() *samurai.FrequencyTable {
	return g.freqTable
}
