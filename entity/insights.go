package entity

type Insight struct {
	ID               string
	ProjectRef       string
	Package          string
	TotalIdentifiers int
	TotalExported    int
	TotalSplits      int
	TotalExpansions  int
	TotalWeight      float64
	Files            []string
}

func (i Insight) AvgSplits() float64 {
	return float64(i.TotalSplits) / float64(i.TotalIdentifiers)
}

func (i Insight) AvgExpansions() float64 {
	return float64(i.TotalExpansions) / float64(i.TotalIdentifiers)
}

func (i Insight) Rate() float64 {
	return i.TotalWeight / float64(i.TotalIdentifiers)
}
