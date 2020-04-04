package usecase

import (
	"strings"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestSplit_OnClosedChannel_ShouldSendNoElements(t *testing.T) {
	identc := make(chan code.Identifier)
	close(identc)

	splitc := split(identc, splitter{})

	var identifiers int
	for range splitc {
		identifiers++
	}

	assert.Equal(t, 0, identifiers)
}

func TestSplit_OnEmptySplitter_ShouldSendElementsWithoutSplits(t *testing.T) {
	identc := make(chan code.Identifier)
	go func() {
		identc <- code.Identifier{
			Name:   "main",
			Splits: make(map[string][]string),
		}
		close(identc)
	}()

	splitc := split(identc, []entity.Splitter{}...)

	splits := make([]code.Identifier, 0)
	for ident := range splitc {
		splits = append(splits, ident)
	}

	assert.Equal(t, 1, len(splits))
	assert.Equal(t, 0, len(splits[0].Splits))
}

func TestSplit_OnOneIdentifierAndTwoSplitters_ShouldSendElementsWithTwoSplits(t *testing.T) {
	identc := make(chan code.Identifier)
	go func() {
		identc <- code.Identifier{
			Name:   "star_wars-II",
			Splits: make(map[string][]string),
		}
		close(identc)
	}()

	byHyphen := splitter{
		name: "hyphen",
		sfunc: func(token string) []string {
			return strings.Split(token, "-")
		},
	}

	byUnderscore := splitter{
		name: "underscore",
		sfunc: func(token string) []string {
			return strings.Split(token, "_")
		},
	}

	splitc := split(identc, byHyphen, byUnderscore)

	splitidents := make([]code.Identifier, 0)
	for ident := range splitc {
		splitidents = append(splitidents, ident)
	}

	assert.Equal(t, 1, len(splitidents))

	splits := splitidents[0].Splits
	assert.Equal(t, []string{"star_wars", "II"}, splits["hyphen"])
	assert.Equal(t, []string{"star", "wars-II"}, splits["underscore"])
}

type splitter struct {
	name  string
	sfunc func(string) []string
}

func (s splitter) Name() string {
	if s.name != "" {
		return s.name
	}

	return "test"
}

func (s splitter) Split(token string) []string {
	if s.sfunc != nil {
		return s.sfunc(token)
	}

	return []string{}
}
