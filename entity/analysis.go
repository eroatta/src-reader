package entity

import (
	"go/ast"
	"go/token"
)

const (
	MinerWordCount          MinerType = "WordCount"
	MinerScopedDeclarations MinerType = "Scoped Declarations"
	MinerDeclarations       MinerType = "Declarations"
)

// AnalysisConfig defines the configuration options for an analysis execution.
type AnalysisConfig struct {
	Miners                    []Miner
	ExtractorFactory          ExtractorFactory
	SplittingAlgorithmFactory SplitterAbstractFactory
	Splitters                 []string
	ExpansionAlgorithmFactory ExpanderAbstractFactory
	Expanders                 []string
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
	Identifiers() []Identifier
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
	Expand(ident Identifier) []string
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

// File TODO
type File struct {
	Name    string
	Raw     []byte
	AST     *ast.File
	FileSet *token.FileSet
	Error   error
}

// Identifier TODO
type Identifier struct {
	File       string
	Position   token.Pos
	Name       string
	Type       string
	Parent     string
	ParentPos  token.Pos
	Node       *ast.Node
	Splits     map[string][]string
	Expansions map[string][]string
	Error      error
}

// ScopedDecl represents the related scope for a declaration.
type ScopedDecl struct {
	ID              string
	DeclType        token.Token
	Name            string
	VariableDecls   []string
	Statements      []string
	BodyText        []string
	Comments        []string
	PackageComments []string
}

// Decl contains the mined text (words and phrases) related to a declaration.
type Decl struct {
	ID       string
	DeclType token.Token
	Words    map[string]struct{}
	Phrases  map[string]struct{}
}
