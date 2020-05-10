package expander

import (
	"errors"
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/port/outgoing/adapter/algorithm/miner"
	"github.com/eroatta/token/basic"
	"github.com/eroatta/token/expansion"
)

// NewBasicFactory creates a new Basic expanders factory.
func NewBasicFactory() entity.ExpanderFactory {
	return basicFactory{}
}

type basicFactory struct{}

func (f basicFactory) Make(miningResults map[string]entity.Miner) (entity.Expander, error) {
	declarationsMiner, ok := miningResults["declarations"]
	if !ok {
		return nil, errors.New("unable to retrieve input from declarations miner")
	}
	declarations := declarationsMiner.Results().(map[string]miner.Decl)

	return &basicExpander{
		expander:     expander{"basic"},
		declarations: declarations,
	}, nil
}

type basicExpander struct {
	expander
	declarations map[string]miner.Decl
}

// Expand receives a entity.Identifier and processes the available splits that
// can be expanded with the current algorithm.
// On Basic, we rely on the related declaration information for the identifier.
// If no declaration information can be found, we avoid trying to expand the identifier
// because results can be broad.
// If a declaration is found but several expansions are found, we handle a subset of them.
func (b basicExpander) Expand(ident entity.Identifier) []entity.Expansion {
	splits, ok := ident.Splits[b.ApplicableOn()]
	if !ok {
		return []entity.Expansion{}
	}

	declarationID := ident.ID
	if ident.IsLocal() {
		declarationID = ident.Parent
	}
	decl, ok := b.declarations[declarationID]
	if !ok {
		expansions := make([]entity.Expansion, len(splits))
		for i, split := range splits {
			expansions[i] = entity.Expansion{From: split.Value, Values: []string{split.Value}}
		}
		return expansions
	}

	wordsBuilder := expansion.NewSetBuilder()
	for k := range decl.Words {
		wordsBuilder.AddStrings(k)
	}
	words := wordsBuilder.Build()

	phrases := make(map[string]string)
	for phrase := range decl.Phrases {
		var acron strings.Builder
		for _, word := range strings.Split(phrase, " ") {
			acron.WriteByte(word[0])
		}

		phrases[acron.String()] = strings.ReplaceAll(phrase, " ", "-")
	}

	expanded := make([]entity.Expansion, len(splits))
	for i, split := range splits {
		expansions := basic.Expand(split.Value, words, phrases, basic.DefaultExpansions)
		if len(expansions) == 0 {
			expansions = []string{split.Value}
		}

		if len(expansions) > 1 {
			expansions = handleMultipleExpansions(split.Value, expansions)
		}

		expanded[i] = entity.Expansion{From: split.Value, Values: expansions}
	}

	return expanded
}

func (b basicExpander) ApplicableOn() string {
	return "greedy"
}

// handleMultipleExpansinos measures the distance between two strings according the
// Levenshtein algorithm, and select the closest three expansions.
func handleMultipleExpansions(token string, expansions []string) []string {
	distances := make([]distance, len(expansions))
	for i, expansion := range expansions {
		value := levenshtein.ComputeDistance(token, expansion)
		distances[i] = distance{value, expansion}
	}
	sort.Sort(byValue(distances))

	limit := 3
	if len(distances) < limit {
		limit = len(distances)
	}

	picked := []string{}
	for _, distance := range distances[0:limit] {
		picked = append(picked, distance.expansion)
	}

	return picked
}

type distance struct {
	value     int
	expansion string
}

type byValue []distance

func (a byValue) Len() int           { return len(a) }
func (a byValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byValue) Less(i, j int) bool { return a[i].value < a[j].value }
