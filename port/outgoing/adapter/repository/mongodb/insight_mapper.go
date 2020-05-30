package mongodb

import (
	"math"

	"github.com/eroatta/src-reader/entity"
)

// insightMapper maps an entity.Insight between its model and database representations.
type insightMapper struct{}

// toDTO maps the entity for entity.Insight into a Data Transfer Object.
func (is *insightMapper) toDTO(ent entity.Insight) insightDTO {
	files := make([]string, 0)
	for file := range ent.Files {
		files = append(files, file)
	}

	avgSplits := make(map[string]float64)
	for alg := range ent.TotalSplits {
		avgSplits[alg] = math.Round(ent.AvgSplits(alg))
	}

	avgExpansions := make(map[string]float64)
	for alg := range ent.TotalExpansions {
		avgExpansions[alg] = math.Round(ent.AvgExpansions(alg))
	}

	return insightDTO{
		ProjectRef:       ent.ProjectRef,
		Package:          ent.Package,
		Accuracy:         ent.Rate(),
		TotalIdentifiers: ent.TotalIdentifiers,
		TotalExported:    ent.TotalExported,
		TotalSplits:      ent.TotalSplits,
		AvgSplits:        avgSplits,
		TotalExpansions:  ent.TotalExpansions,
		AvgExpansions:    avgExpansions,
		TotalWeight:      ent.TotalWeight,
		Files:            files,
	}
}

type insightDTO struct {
	ID               string             `bson:"_id"`
	ProjectRef       string             `bson:"project_ref"`
	Package          string             `bson:"package"`
	Accuracy         float64            `bson:"accuracy"`
	TotalIdentifiers int                `bson:"total_identifiers"`
	TotalExported    int                `bson:"total_exported"`
	TotalSplits      map[string]int     `bson:"total_splits"`
	AvgSplits        map[string]float64 `bson:"avg_splits"`
	TotalExpansions  map[string]int     `bson:"total_expansions"`
	AvgExpansions    map[string]float64 `bson:"avg_expansions"`
	TotalWeight      float64            `bson:"total_weight"`
	Files            []string           `bson:"files"`
}
