package entity

import (
	"go/ast"

	"github.com/eroatta/src-reader/code"
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
	StaticInputs     map[string]lists.List
	Miners           []Miner
	ExtractorFactory ExtractorFactory
	Splitters        []Splitter
	Expanders        []Expander
}

type MinerType string

// Miner interface is used to define a custom miner.
type Miner interface {
	// Type provides the name of the miner.
	Type() MinerType
	// Visit applies the mining logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
}

type ExtractorFactory func(filename string) Extractor

// Extractor is used to define a custom identifier extractor.
type Extractor interface {
	// Visit applies the extraction logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
	// Identifiers returns the extracted identifiers.
	Identifiers() []code.Identifier
}

// Splitter interface is used to define a custom splitter.
type Splitter interface {
	// Name returns the name of the custom splitter.
	Name() string
	// Split returns the split identifier.
	Split(token string) []string //TODO: we should return a string, divided by something like hyphen
}

// Expander interface is used to define a custom expander.
type Expander interface {
	// Name returns the name of the custom expander.
	Name() string
	// ApplicableOn defines the name of splits used as input.
	ApplicableOn() string
	// Expand performs the expansion on the token as a whole.
	Expand(ident code.Identifier) []string
}