package mongodb

import (
	"testing"
	"time"

	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestToDTO_OnAnalysisMapper_ShouldReturnAnalysisDTO(t *testing.T) {
	now := time.Now()
	ent := entity.AnalysisResults{
		ID:                      "715f17550be5f7222a815ff80966adaf",
		ProjectName:             "src-d/go-siva",
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

	assert.Equal(t, "715f17550be5f7222a815ff80966adaf", dto.ID)
	assert.Equal(t, "src-d/go-siva", dto.ProjectRef)
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
