package splitter

import (
	"github.com/eroatta/src-reader/step"
	"github.com/eroatta/token/conserv"
)

type conservSplitter struct {
	splitter
}

// Split splits a token using the Conserv splitter.
func (c conservSplitter) Split(token string) []string {
	return conserv.Split(token)
}

// NewConserv creates a new Conserv splitter.
func NewConserv() step.Splitter {
	return conservSplitter{
		splitter: splitter{"conserv"},
	}
}
