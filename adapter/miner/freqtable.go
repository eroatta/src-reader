package miner

import (
	"go/ast"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/samurai"
)

// GlobalFreqTable represents a Global Frequency Table miner.
type GlobalFreqTable struct {
	freqTable *samurai.FrequencyTable
}

// NewGlobalFreqTable creates a new GlobalFreqTable miner with the provided frequency table.
func NewGlobalFreqTable(ft *samurai.FrequencyTable) GlobalFreqTable {
	return GlobalFreqTable{
		freqTable: ft,
	}
}

// Type specifies the GlobalFreqTable miner type.
func (g GlobalFreqTable) Type() entity.MinerType {
	return entity.MinerGlobalFrequencyTable
}

// SetCurrentFile changes the current file on the miner using the provided filename.
func (g GlobalFreqTable) SetCurrentFile(filename string) {
	// do nothing
}

// Visit traverses a ast.Node for mining.
func (g GlobalFreqTable) Visit(node ast.Node) ast.Visitor {
	return nil
}

// Table returns the frequency table after the mining process.
func (g GlobalFreqTable) Table() *samurai.FrequencyTable {
	return g.freqTable
}
