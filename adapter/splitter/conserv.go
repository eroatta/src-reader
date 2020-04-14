package splitter

import (
	"strings"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/conserv"
)

// NewConservFactory creates a new Conserv splitter factory.
func NewConservFactory() entity.SplitterFactory {
	return conservFactory{}
}

type conservFactory struct{}

func (f conservFactory) Make(miningResults map[entity.MinerType]entity.Miner) (entity.Splitter, error) {
	return conservSplitter{
		splitter: splitter{"conserv"},
	}, nil
}

type conservSplitter struct {
	splitter
}

// Split splits a token using the Conserv splitter.
func (c conservSplitter) Split(token string) []entity.Split {
	splits := []entity.Split{}
	for i, split := range strings.Split(conserv.Split(token), " ") {
		splits = append(splits, entity.Split{Order: i + 1, Value: strings.ToLower(split)})
	}

	return splits
}
