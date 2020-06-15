package mongodb

import (
	"testing"
	"time"

	"github.com/eroatta/src-reader/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToDTO_OnAnalysisMapper_ShouldReturnAnalysisDTO(t *testing.T) {
	now := time.Now()
	ent := entity.AnalysisResults{
		ID:                      uuid.MustParse("f9b76fde-c342-4328-8650-85da8f21e2be"),
		ProjectName:             "src-d/go-siva",
		ProjectID:               uuid.MustParse("f9b76fde-c342-4328-8650-85da8f21e2be"),
		DateCreated:             now,
		PipelineMiners:          []string{"miner_1", "miner_2"},
		PipelineSplitters:       []string{"splitter_1", "splitter_2"},
		PipelineExpanders:       []string{"expander_1", "expander_2"},
		FilesTotal:              10,
		FilesValid:              8,
		FilesError:              2,
		FilesErrorSamples:       []string{"file_error"},
		IdentifiersTotal:        120,
		IdentifiersValid:        105,
		IdentifiersError:        15,
		IdentifiersErrorSamples: []string{"identifier_error"},
	}

	am := &analysisMapper{}
	dto := am.toDTO(ent)

	assert.Equal(t, "f9b76fde-c342-4328-8650-85da8f21e2be", dto.ID)
	assert.Equal(t, "src-d/go-siva", dto.ProjectRef)
	assert.Equal(t, "f9b76fde-c342-4328-8650-85da8f21e2be", dto.ProjectID)
	assert.Equal(t, now, dto.CreatedAt)
	assert.ElementsMatch(t, []string{"miner_1", "miner_2"}, dto.Miners)
	assert.ElementsMatch(t, []string{"splitter_1", "splitter_2"}, dto.Splitters)
	assert.ElementsMatch(t, []string{"expander_1", "expander_2"}, dto.Expanders)
	assert.Equal(t, int32(10), dto.Files.Total)
	assert.Equal(t, int32(8), dto.Files.Valid)
	assert.Equal(t, int32(2), dto.Files.Failed)
	assert.ElementsMatch(t, []string{"file_error"}, dto.Files.ErrorSamples)
	assert.Equal(t, int32(120), dto.Identifiers.Total)
	assert.Equal(t, int32(105), dto.Identifiers.Valid)
	assert.Equal(t, int32(15), dto.Identifiers.Failed)
	assert.ElementsMatch(t, []string{"identifier_error"}, dto.Identifiers.ErrorSamples)
}

func TestToEntity_OnAnalysisMapper_ShouldReturnAnalysisResultsEntity(t *testing.T) {
	now := time.Now()
	dto := analysisDTO{
		ID:         "f9b76fde-c342-4328-8650-85da8f21e2be",
		CreatedAt:  now,
		ProjectRef: "src-d/go-siva",
		ProjectID:  "f9b76fde-c342-4328-8650-85da8f21e2be",
		Miners:     []string{"miner_1", "miner_2"},
		Splitters:  []string{"splitter_1", "splitter_2"},
		Expanders:  []string{"expander_1", "expander_2"},
		Files: summarizerDTO{
			Total:        10,
			Valid:        8,
			Failed:       2,
			ErrorSamples: []string{"file_error"},
		},
		Identifiers: summarizerDTO{
			Total:        120,
			Valid:        105,
			Failed:       15,
			ErrorSamples: []string{"identifier_error"},
		},
	}

	am := &analysisMapper{}
	ent := am.toEntity(dto)

	assert.Equal(t, uuid.MustParse("f9b76fde-c342-4328-8650-85da8f21e2be"), ent.ID)
	assert.Equal(t, "src-d/go-siva", ent.ProjectName)
	assert.Equal(t, uuid.MustParse("f9b76fde-c342-4328-8650-85da8f21e2be"), ent.ProjectID)
	assert.Equal(t, now, ent.DateCreated)
	assert.ElementsMatch(t, []string{"miner_1", "miner_2"}, ent.PipelineMiners)
	assert.ElementsMatch(t, []string{"splitter_1", "splitter_2"}, ent.PipelineSplitters)
	assert.ElementsMatch(t, []string{"expander_1", "expander_2"}, ent.PipelineExpanders)
	assert.Equal(t, 10, ent.FilesTotal)
	assert.Equal(t, 8, ent.FilesValid)
	assert.Equal(t, 2, ent.FilesError)
	assert.ElementsMatch(t, []string{"file_error"}, ent.FilesErrorSamples)
	assert.Equal(t, 120, ent.IdentifiersTotal)
	assert.Equal(t, 105, ent.IdentifiersValid)
	assert.Equal(t, 15, ent.IdentifiersError)
	assert.ElementsMatch(t, []string{"identifier_error"}, ent.IdentifiersErrorSamples)
}
