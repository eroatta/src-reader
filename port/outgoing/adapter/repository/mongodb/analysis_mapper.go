package mongodb

import (
	"time"

	"github.com/eroatta/src-reader/entity"
)

// analysisMapper maps an AnalysisResults entity between its model and database representation.
type analysisMapper struct{}

// toDTO maps the entity for AnalysisResults into a Data Transfer Object.
func (am *analysisMapper) toDTO(ent entity.AnalysisResults) analysisDTO {
	return analysisDTO{
		ID:         ent.ID,
		CreatedAt:  ent.DateCreated,
		ProjectRef: ent.ProjectName,
		Miners:     ent.PipelineMiners,
		Splitters:  ent.PipelineSplitters,
		Expanders:  ent.PipelineExpanders,
		Files: summarizerDTO{
			Total:        int32(ent.FilesTotal),
			Valid:        int32(ent.FilesValid),
			Failed:       int32(ent.FilesError),
			ErrorSamples: ent.FilesErrorSamples,
		},
		Identifiers: summarizerDTO{
			Total:        int32(ent.IdentifiersTotal),
			Valid:        int32(ent.IdentifiersValid),
			Failed:       int32(ent.IdentifiersError),
			ErrorSamples: ent.IdentifiersErrorSamples,
		},
	}
}

// analysisDTO is the database representation for an AnalysisResults.
type analysisDTO struct {
	ID          string        `bson:"_id"`
	CreatedAt   time.Time     `bson:"created_at"`
	ProjectRef  string        `bson:"project_ref"`
	Miners      []string      `bson:"miners"`
	Splitters   []string      `bson:"splitters"`
	Expanders   []string      `bson:"expanders"`
	Files       summarizerDTO `bson:"files_summary"`
	Identifiers summarizerDTO `bson:"identifiers_summary"`
}

// summarizerDTO is the database representation for an AnalysisResults summary.
type summarizerDTO struct {
	Total        int32    `bson:"total"`
	Valid        int32    `bson:"valid"`
	Failed       int32    `bson:"failed"`
	ErrorSamples []string `bson:"error_samples"`
}
