package entity

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/agnivade/levenshtein"
	"github.com/google/uuid"
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
	ID            string
	Package       string
	File          string
	Position      token.Pos
	Name          string
	Type          token.Token
	Node          *ast.Node
	Splits        map[string][]Split
	Expansions    map[string][]Expansion
	Error         error
	Normalization Normalization
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

// Exported determines if the identifier is exported on its package.
func (i Identifier) Exported() bool {
	if len(i.Name) == 0 {
		return false
	}

	return unicode.IsUpper(rune(i.Name[0])) && unicode.IsLetter(rune(i.Name[0]))
}

// Normalize applies a normalization function to select the best split/expansion approach.
func (i *Identifier) Normalize() {
	normalization := Normalization{
		Word:      "undefined",
		Algorithm: "undefined",
		Score:     0.0,
	}

	for algorithm, expansions := range i.Expansions {
		sort.Sort(bySoftwordOrder(expansions))

		var wordBuilder strings.Builder
		for i, expansion := range expansions {
			// TODO: review
			if len(expansion.Values) > 0 {
				expandedSoftword := expansion.Values[0] // TODO: for now, pick the first one
				if i > 0 {
					expandedSoftword = strings.Title(expandedSoftword)
				}
				wordBuilder.WriteString(expandedSoftword)
			}
		}
		word := wordBuilder.String()

		lengths := float64(len(i.Name) + len(word))
		score := (lengths - float64(levenshtein.ComputeDistance(i.Name, word))) / lengths
		if score >= normalization.Score {
			normalization = Normalization{
				Word:      word,
				Algorithm: fmt.Sprintf("%s+%s", expansions[0].SplittingAlgorithm, algorithm),
				Score:     score,
			}
		}
	}

	i.Normalization = normalization
}

// Split represents a hardword or softword in which the identifier was divided.
type Split struct {
	Order int
	Value string
}

// Expansion represents a set of expansions from a split.
type Expansion struct {
	Order              int
	From               string
	Values             []string
	SplittingAlgorithm string
}

type bySoftwordOrder []Expansion

func (a bySoftwordOrder) Len() int           { return len(a) }
func (a bySoftwordOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bySoftwordOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

// Normalization represents a word composition with its given score.
type Normalization struct {
	Word      string
	Algorithm string
	Score     float64
}

// AnalysisResults represents the results for an analysis, indicating its creation date,
// the configuration provided (URL, miners, splitters, expanders), and information about
// the processed files and identifiers.
type AnalysisResults struct {
	ID                      uuid.UUID
	DateCreated             time.Time
	ProjectID               uuid.UUID
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
