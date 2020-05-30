package mongodb

import (
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestToDTO_OnInsightMapper_ShouldReturnInsightDTO(t *testing.T) {
	ent := entity.Insight{

		ProjectRef:       "eroatta/test",
		Package:          "main",
		TotalIdentifiers: 3,
		TotalExported:    1,
		TotalSplits: map[string]int{
			"conserv": 5,
		},
		TotalExpansions: map[string]int{
			"no_exp": 5,
		},
		TotalWeight: 2.267,
		Files: map[string]struct{}{
			"main.go":   {},
			"helper.go": {},
		},
	}

	im := &insightMapper{}
	dto := im.toDTO(ent)

	assert.Equal(t, "", dto.ID)
	assert.Equal(t, "eroatta/test", dto.ProjectRef)
	assert.Equal(t, "main", dto.Package)
	assert.Equal(t, ent.Rate(), dto.Accuracy)
	assert.Equal(t, 3, dto.TotalIdentifiers)
	assert.Equal(t, 1, dto.TotalExported)
	assert.EqualValues(t, map[string]int{
		"conserv": 5,
	}, dto.TotalSplits)
	assert.EqualValues(t, map[string]float64{
		"conserv": 2.0,
	}, dto.AvgSplits)
	assert.EqualValues(t, map[string]int{
		"no_exp": 5,
	}, dto.TotalExpansions)
	assert.EqualValues(t, map[string]float64{
		"no_exp": 2.0,
	}, dto.AvgExpansions)
	assert.Equal(t, 2.267, dto.TotalWeight)
	assert.ElementsMatch(t, []string{"main.go", "helper.go"}, dto.Files)
}
