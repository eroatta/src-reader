package splitter

import (
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/conserv"
	"github.com/eroatta/token/lists"
)

// NewConservFactory creates a new Conserv splitter factory.
func NewConservFactory() entity.SplitterFactory {
	return conservFactory{}
}

type conservFactory struct{}

func (f conservFactory) Make(staticInputs map[string]lists.List, miningResults map[entity.MinerType]entity.Miner) (entity.Splitter, error) {
	return conservSplitter{
		splitter: splitter{"conserv"},
	}, nil
}

type conservSplitter struct {
	splitter
}

// Split splits a token using the Conserv splitter.
func (c conservSplitter) Split(token string) []string {
	return conserv.Split(token)
}
