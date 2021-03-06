package step_test

import (
	"strings"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/usecase/step"
	"github.com/stretchr/testify/assert"
)

func TestSplit_OnClosedChannel_ShouldSendNoElements(t *testing.T) {
	identc := make(chan entity.Identifier)
	close(identc)

	splitc := step.Split(identc, splitter{})

	var identifiers int
	for range splitc {
		identifiers++
	}

	assert.Equal(t, 0, identifiers)
}

func TestSplit_OnEmptySplitter_ShouldSendElementsWithoutSplits(t *testing.T) {
	identc := make(chan entity.Identifier)
	go func() {
		identc <- entity.Identifier{
			Name:   "main",
			Splits: make(map[string][]entity.Split),
		}
		close(identc)
	}()

	splitc := step.Split(identc, []entity.Splitter{}...)

	splits := make([]entity.Identifier, 0)
	for ident := range splitc {
		splits = append(splits, ident)
	}

	assert.Equal(t, 1, len(splits))
	assert.Equal(t, 0, len(splits[0].Splits))
}

func TestSplit_OnOneIdentifierAndTwoSplitters_ShouldSendElementsWithTwoSplits(t *testing.T) {
	identc := make(chan entity.Identifier)
	go func() {
		identc <- entity.Identifier{
			Name:   "star_wars-II",
			Splits: make(map[string][]entity.Split),
		}
		close(identc)
	}()

	byHyphen := splitter{
		name: "hyphen",
		sfunc: func(token string) string {
			return strings.ReplaceAll(token, "-", " ")
		},
	}

	byUnderscore := splitter{
		name: "underscore",
		sfunc: func(token string) string {
			return strings.ReplaceAll(token, "_", " ")
		},
	}

	splitc := step.Split(identc, byHyphen, byUnderscore)

	splitidents := make([]entity.Identifier, 0)
	for ident := range splitc {
		splitidents = append(splitidents, ident)
	}

	assert.Equal(t, 1, len(splitidents))

	splits := splitidents[0].Splits
	assert.Equal(t, []entity.Split{{Order: 1, Value: "star_wars"}, {Order: 2, Value: "II"}}, splits["hyphen"])
	assert.Equal(t, []entity.Split{{Order: 1, Value: "star"}, {Order: 2, Value: "wars-II"}}, splits["underscore"])
}

type splitter struct {
	name  string
	sfunc func(string) string
}

func (s splitter) Name() string {
	if s.name != "" {
		return s.name
	}

	return "test"
}

func (s splitter) Split(token string) []entity.Split {
	if s.sfunc != nil {
		splits := []entity.Split{}
		for i, sp := range strings.Split(s.sfunc(token), " ") {
			splits = append(splits, entity.Split{Order: i + 1, Value: sp})
		}

		return splits
	}

	return []entity.Split{}
}
