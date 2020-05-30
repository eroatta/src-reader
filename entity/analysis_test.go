package entity_test

import (
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestFullPackageName_OnIdentifier(t *testing.T) {
	cases := []struct {
		ident    entity.Identifier
		expected string
	}{
		{ident: entity.Identifier{Package: "main", File: "main.go"}, expected: "main"},
		{ident: entity.Identifier{Package: "main", File: "helper.go"}, expected: "main"},
		{ident: entity.Identifier{Package: "amap", File: "amap/pattern.go"}, expected: "amap"},
		{ident: entity.Identifier{Package: "persistence", File: "adapter/persistence/postgresql.go"}, expected: "adapter/persistence"},
		{ident: entity.Identifier{Package: "crypto", File: "utils/crypto.go"}, expected: "utils/crypto"},
		{ident: entity.Identifier{}, expected: ""},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.ident.FullPackageName())
	}
}

func TestExported_OnIdentifier(t *testing.T) {
	cases := []struct {
		ident    entity.Identifier
		expected bool
	}{
		{ident: entity.Identifier{Name: "helperFunc"}, expected: false},
		{ident: entity.Identifier{Name: "HelperFunc"}, expected: true},
		{ident: entity.Identifier{Name: ""}, expected: false},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.ident.Exported())
	}
}

func TestNormalize_OnIdentifier(t *testing.T) {
	cases := []struct {
		name     string
		ident    entity.Identifier
		expected entity.Normalization
	}{
		{
			name: "no_expansions",
			ident: entity.Identifier{
				Expansions: make(map[string][]entity.Expansion),
			},
			expected: entity.Normalization{
				Word:      "undefined",
				Algorithm: "undefined",
				Score:     0.0,
			},
		},
		{
			name: "one_expansion_perfect_score",
			ident: entity.Identifier{
				Name: "delete",
				Expansions: map[string][]entity.Expansion{
					"custom_test": {
						{Order: 1, SpittingAlgorithm: "split_test", From: "delete", Values: []string{"delete"}},
					},
				},
			},
			expected: entity.Normalization{
				Word:      "delete",
				Algorithm: "split_test+custom_test",
				Score:     1.0,
			},
		},
		{
			name: "one_expansion_close_score",
			ident: entity.Identifier{
				Name: "ctrldel",
				Expansions: map[string][]entity.Expansion{
					"custom": {
						{Order: 1, SpittingAlgorithm: "custom", From: "ctrl", Values: []string{"control"}},
						{Order: 2, SpittingAlgorithm: "custom", From: "del", Values: []string{"delete"}},
					},
				},
			},
			expected: entity.Normalization{
				Word:      "controlDelete",
				Algorithm: "custom+custom",
				Score:     0.65,
			},
		},
		{
			name: "two_expansions_close_score",
			ident: entity.Identifier{
				Name: "ctrldel",
				Expansions: map[string][]entity.Expansion{
					"custom": {
						{Order: 1, SpittingAlgorithm: "custom", From: "ctrl", Values: []string{"control"}},
						{Order: 2, SpittingAlgorithm: "custom", From: "del", Values: []string{"delete"}},
					},
					"another_custom": {
						{Order: 1, SpittingAlgorithm: "custom", From: "ctrl", Values: []string{"control"}},
						{Order: 2, SpittingAlgorithm: "custom", From: "del", Values: []string{"delay", "delete"}},
					},
				},
			},
			expected: entity.Normalization{
				Word:      "controlDelay",
				Algorithm: "custom+another_custom",
				Score:     0.6842105263157895,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.ident.Normalize()

			assert.Equal(t, c.expected, c.ident.Normalization)
		})
	}
}
