package splitter

import (
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/greedy"
)

type greedySplitter struct {
	splitter
}

// Split splits a token using the Greedy splitter.
func (g greedySplitter) Split(token string) []string {
	return greedy.Split(token, greedy.DefaultList)
}

// NewGreedy creates a new Greedy splitter.
func NewGreedy() entity.Splitter {
	return greedySplitter{
		splitter: splitter{"greedy"},
	}
}
