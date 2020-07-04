package entity

import "github.com/google/uuid"

// Insight represents information extracted and summarized from an Analysis, for a package.
type Insight struct {
	ID               string
	ProjectRef       string
	AnalysisID       uuid.UUID
	Package          string
	TotalIdentifiers int
	TotalExported    int
	TotalSplits      map[string]int
	TotalExpansions  map[string]int
	TotalWeight      float64
	Files            map[string]struct{}
}

// AvgSplits returns the average number of splits for a particular algorithm.
func (i Insight) AvgSplits(algorithm string) float64 {
	return float64(i.TotalSplits[algorithm]) / float64(i.TotalIdentifiers)
}

// AvgExpansions returns the average number of expansions for a particular algorithm.
func (i Insight) AvgExpansions(algorithm string) float64 {
	return float64(i.TotalExpansions[algorithm]) / float64(i.TotalIdentifiers)
}

// Rate returns the correctness rate for the Insight.
func (i Insight) Rate() float64 {
	return i.TotalWeight / float64(i.TotalIdentifiers)
}

// Include includes an identifier into the analysis for the current package Insight.
func (i *Insight) Include(ident Identifier) {
	if i.Package != ident.FullPackageName() {
		return
	}

	i.TotalIdentifiers++
	if ident.Exported() {
		i.TotalExported++
	}

	for algorithm, splits := range ident.Splits {
		i.TotalSplits[algorithm] += len(splits)
	}

	for algorithm, expansions := range ident.Expansions {
		i.TotalExpansions[algorithm] += len(expansions)
	}

	weight := 0.7
	if ident.Exported() {
		weight = 1.0
	}
	i.TotalWeight += ident.Normalization.Score * weight

	i.Files[ident.File] = struct{}{}
}
