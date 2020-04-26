package entity

import (
	"go/ast"
	"go/token"
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

type Split struct {
	Order int
	Value string
}

type Expansion struct {
	From   string
	Values []string
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
