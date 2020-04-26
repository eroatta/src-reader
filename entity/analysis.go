package entity

import (
	"go/ast"
	"go/token"
	"time"
)

// AnalysisConfig defines the configuration options for an analysis execution.
type AnalysisConfig struct {
	Miners                    []string
	MinerAlgorithmFactory     MinerAbstractFactory
	ExtractorFactory          ExtractorFactory
	SplittingAlgorithmFactory SplitterAbstractFactory
	Splitters                 []string
	ExpansionAlgorithmFactory ExpanderAbstractFactory
	Expanders                 []string
}

// File represents a source code file, including its raw form and also its Abstract Syntax Tree representation.
type File struct {
	Name    string
	Raw     []byte
	AST     *ast.File
	FileSet *token.FileSet
	Error   error
}

// Identifier represents an identifier extracted from source code, indicating its origin, type,
// parent information, and splits/expansions.
type Identifier struct {
	ID         string
	File       string
	Position   token.Pos
	Name       string
	Type       token.Token
	Parent     string
	ParentPos  token.Pos
	Node       *ast.Node
	Splits     map[string][]Split
	Expansions map[string][]Expansion
	Error      error
}

// IsLocal indicates if an identifier is part of a function.
func (i Identifier) IsLocal() bool {
	return i.Parent != ""
}

// Split represents a hardword or softword in which the identifier was divided.
type Split struct {
	Order int
	Value string
}

// Expansion represents a set of expansions from a split.
type Expansion struct {
	From   string
	Values []string
}

// AnalysisResults represents the results for an analysis, indicating its creation date,
// the configuration provided (URL, miners, splitters, expanders), and information about
// the processed files and identifiers.
type AnalysisResults struct {
	// id
	// status (?)
	DateCreated             time.Time
	ProjectID               string
	ProjectURL              string
	PipelineMiners          []string
	PipelineSplitters       []string
	PipelineExpanders       []string
	FilesTotal              int
	FilesValid              int
	FilesError              int
	FilesErrorSamples       []string
	IdentifiersTotal        int
	IdentifiersValid        int
	IdentifiersError        int
	IdentifiersErrorSamples []string
}
