package entity

import "go/ast"

const (
	MinerWordCount            MinerType = "WordCount"
	MinerScopedDeclarations   MinerType = "Scoped Declarations"
	MinerDeclarations         MinerType = "Declarations"
	MinerComments             MinerType = "Comments"
	MinerGlobalFrequencyTable MinerType = "Global Frequency Table"
)

type MinerType string

// Miner interface is used to define a custom miner.
type Miner interface {
	// Type provides the name of the miner.
	Type() MinerType
	// Visit applies the mining logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
	// SetCurrentFile specifies the current file being mined.
	SetCurrentFile(filename string)
}

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

// Expander interface is used to define a custom expander.
type Expander interface {
	// Name returns the name of the custom expander.
	Name() string
	// ApplicableOn defines the name of splits used as input.
	ApplicableOn() string
	// Expand performs the expansion on the token as a whole.
	Expand(ident Identifier) []Expansion
}

type SplitterAbstractFactory interface {
	Get(algorithm string) (SplitterFactory, error)
}

type SplitterFactory interface {
	Make(miningResults map[MinerType]Miner) (Splitter, error)
}

type ExpanderAbstractFactory interface {
	Get(algorithm string) (ExpanderFactory, error)
}

type ExpanderFactory interface {
	Make(miningResults map[MinerType]Miner) (Expander, error)
}
