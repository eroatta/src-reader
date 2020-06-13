package step_test

import (
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/usecase/step"
	"github.com/stretchr/testify/assert"
)

func TestNormalize_OnClosedChannel_ShouldSendNoElements(t *testing.T) {
	identc := make(chan entity.Identifier)
	close(identc)

	normalizedc := step.Normalize(identc)

	var identifiers int
	for range normalizedc {
		identifiers++
	}

	assert.Equal(t, 0, identifiers)
}

func TestNormalize_OnIdentifier_ShouldSendNormalizedElement(t *testing.T) {
	identc := make(chan entity.Identifier)
	go func() {
		identc <- entity.Identifier{
			Name: "ctrldel",
			Splits: map[string][]entity.Split{
				"custom": {
					{Order: 1, Value: "ctrl"},
					{Order: 2, Value: "del"},
				},
			},
			Expansions: map[string][]entity.Expansion{
				"custom": {
					{Order: 1, SplittingAlgorithm: "custom", From: "ctrl", Values: []string{"control"}},
					{Order: 2, SplittingAlgorithm: "custom", From: "del", Values: []string{"delete"}},
				},
			},
		}
		close(identc)
	}()

	normalizedc := step.Normalize(identc)

	normalized := make([]entity.Identifier, 0)
	for ident := range normalizedc {
		normalized = append(normalized, ident)
	}

	assert.Equal(t, 1, len(normalized))

	ident := normalized[0]
	assert.Equal(t, "controlDelete", ident.Normalization.Word)
	assert.Equal(t, "custom+custom", ident.Normalization.Algorithm)
	assert.Equal(t, 0.65, ident.Normalization.Score)
}
