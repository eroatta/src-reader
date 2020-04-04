package entity

import (
	"go/ast"

	"github.com/eroatta/token/lists"
)

const (
	Words                MinerType = "Words"
	Phrases              MinerType = "Phrases"
	LocalFrequencyTable  MinerType = "Local Frequency Table"
	GlobalFrequencyTable MinerType = "Global Frequency Table"
	ScopedDeclarations   MinerType = "Scoped Declarations"
)

// AnalysisConfig defines the configuration options for an analysis execution.
type AnalysisConfig struct {
	StaticInputs map[string]lists.List
	Miners       []Miner
	Splitters    []Splitter
	Expanders    []Expander
}

type MinerType string

// Miner interface is used to define a custom miner.
type Miner interface {
	// Type provides the name of the miner.
	Type() MinerType
	// Visit applies the mining logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
}

type Splitter interface {
}

type Expander interface {
}
