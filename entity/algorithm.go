package entity

import "go/ast"

// Miner interface is used to define a custom miner.
type Miner interface {
	// Name provides the name of the miner.
	Name() string
	// Visit applies the mining logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
	// SetCurrentFile specifies the current file being mined.
	SetCurrentFile(filename string)
	// Results returns the results after mining.
	Results() interface{}
}

// MinerAbstractFactory is an interface for creating mining algorithm factories.
type MinerAbstractFactory interface {
	// Get returns a MinerFactory for the selectd mining algorithm.
	Get(algorithm string) (MinerFactory, error)
}

// MinerFactory is an interface for creating mining algorithm instances.
type MinerFactory interface {
	// Make returns a mining algorithm instance.
	Make() (Miner, error)
}

// TODO: move it?
type ExtractorFactory func(filename string) Extractor

// Extractor is used to define a custom identifier extractor.
type Extractor interface {
	// Visit applies the extraction logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
	// Identifiers returns the extracted identifiers.
	Identifiers() []Identifier
}

// Splitter interface is used to define a custom splitter.
type Splitter interface {
	// Name returns the name of the custom splitter.
	Name() string
	// Split returns the split identifier.
	Split(token string) []Split
}

// SplitterAbstractFactory is an interface for creating splitting algorithm factories.
type SplitterAbstractFactory interface {
	// Get returns a SplitterFactory for the selectd splitting algorithm.
	Get(algorithm string) (SplitterFactory, error)
}

// SplitterFactory is an interface for creating splitting algorithm instances.
type SplitterFactory interface {
	// Make returns a splitting algorithm instance built from miners.
	Make(miners map[string]Miner) (Splitter, error)
}

// Expander interface is used to define a custom expander.
type Expander interface {
	// Name returns the name of the custom expander.
	Name() string
	// ApplicableOn defines the name of splits used as input.
	ApplicableOn() string
	// Expand performs the expansion on the token as a whole.
	Expand(ident Identifier) []Expansion
}

// ExpanderAbstractFactory is an interface for creating expandion algorithm factories.
type ExpanderAbstractFactory interface {
	// Get returns a ExpanderFactory for the selectd expansion algorithm.
	Get(algorithm string) (ExpanderFactory, error)
}

// ExpanderFactory is an interface for creating expansion algorithm instances.
type ExpanderFactory interface {
	// Make returns an expansion algorithm instance built from miners.
	Make(miningResults map[string]Miner) (Expander, error)
}
