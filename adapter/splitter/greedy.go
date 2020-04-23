package splitter

import (
	"strings"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/greedy"
)

// NewGreedyFactory creates a new Greedy splitter factory.
func NewGreedyFactory() entity.SplitterFactory {
	return greedyFactory{}
}

type greedyFactory struct{}

func (f greedyFactory) Make(miningResults map[string]entity.Miner) (entity.Splitter, error) {
	return greedySplitter{
		splitter: splitter{"greedy"},
	}, nil
}

type greedySplitter struct {
	splitter
}

// Split splits a token using the Greedy splitter.
func (g greedySplitter) Split(token string) []entity.Split {
	splits := []entity.Split{}
	for i, split := range strings.Split(greedy.Split(token, greedy.DefaultList), " ") {
		splits = append(splits, entity.Split{Order: i + 1, Value: split})
	}

	return splits
}
