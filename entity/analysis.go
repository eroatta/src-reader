package entity

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
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
	Package    string
	File       string
	Position   token.Pos
	Name       string
	Type       token.Token
	Node       *ast.Node
	Splits     map[string][]Split
	Expansions map[string][]Expansion
	Error      error
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

// FullPackageName returns the package name, including its directory structure.
func (i Identifier) FullPackageName() string {
	idx := strings.LastIndex(i.File, "/")
	if idx == -1 {
		return i.Package
	}

	dir := i.File[:idx]
	if strings.HasSuffix(dir, i.Package) {
		return dir
	}

	return fmt.Sprintf("%s/%s", dir, i.Package)
}

// AnalysisResults represents the results for an analysis, indicating its creation date,
// the configuration provided (URL, miners, splitters, expanders), and information about
// the processed files and identifiers.
type AnalysisResults struct {
	ID                      string
	DateCreated             time.Time
	ProjectName             string
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
