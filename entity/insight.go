package entity

type Insight struct {
	ID               string
	ProjectRef       string
	Package          string
	TotalIdentifiers int
	TotalExported    int
	TotalSplits      map[string]int
	TotalExpansions  map[string]int
	TotalWeight      float64
	Files            map[string]struct{}
}

func (i Insight) AvgSplits(algorithm string) float64 {
	return float64(i.TotalSplits[algorithm]) / float64(i.TotalIdentifiers)
}

func (i Insight) AvgExpansions(algorithm string) float64 {
	return float64(i.TotalExpansions[algorithm]) / float64(i.TotalIdentifiers)
}

func (i Insight) Rate() float64 {
	return i.TotalWeight / float64(i.TotalIdentifiers)
}

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

	//i.TotalWeight += ident.SimilarityRate() * 1 // TODO: adjust weight

	i.Files[ident.File] = struct{}{}
}
