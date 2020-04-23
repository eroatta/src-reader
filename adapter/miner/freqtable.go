package miner

import (
	"go/ast"

	"github.com/eroatta/src-reader/adapter/frequencytable"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/samurai"
)

// NewGlobalFreqTableFactory creates a new declarations miner factory.
func NewGlobalFreqTableFactory() entity.MinerFactory {
	return globalFreqTableFactory{}
}

type globalFreqTableFactory struct{}

func (f globalFreqTableFactory) Make() (entity.Miner, error) {
	return NewGlobalFreqTable(frequencytable.Global), nil
}

// NewGlobalFreqTable creates a new GlobalFreqTable miner with the provided frequency table.
func NewGlobalFreqTable(ft *samurai.FrequencyTable) GlobalFreqTable {
	return GlobalFreqTable{
		miner:     miner{"global-frequency-table"},
		freqTable: ft,
	}
}

// GlobalFreqTable represents a Global Frequency Table miner.
type GlobalFreqTable struct {
	miner
	freqTable *samurai.FrequencyTable
}

// SetCurrentFile changes the current file on the miner using the provided filename.
func (g GlobalFreqTable) SetCurrentFile(filename string) {
	// do nothing
}

// Visit traverses a ast.Node for mining.
func (g GlobalFreqTable) Visit(node ast.Node) ast.Visitor {
	return nil
}

// Results returns the global *samurai.FrequencyTable after the mining process.
func (g GlobalFreqTable) Results() interface{} {
	return g.freqTable
}
