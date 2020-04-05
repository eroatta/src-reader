package splitter

import (
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/greedy"
)

// NewGreedyFactory creates a new Greedy splitter factory.
func NewGreedyFactory() entity.SplitterFactory {
	return greedyFactory{}
}

type greedyFactory struct{}

func (f greedyFactory) Make(staticInputs map[string]interface{}, miningResults map[entity.MinerType]entity.Miner) (entity.Splitter, error) {
	return greedySplitter{
		splitter: splitter{"greedy"},
	}, nil
}

type greedySplitter struct {
	splitter
}

// Split splits a token using the Greedy splitter.
func (g greedySplitter) Split(token string) []string {
	return greedy.Split(token, greedy.DefaultList)
}
