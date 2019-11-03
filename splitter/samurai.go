package splitter

import (
	"github.com/eroatta/src-reader/step"
	"github.com/eroatta/token/lists"
	"github.com/eroatta/token/samurai"
)

type samuraiSplitter struct {
	splitter
	context samurai.TokenContext
}

// Split splits a token using the Samurai splitter.
func (s samuraiSplitter) Split(token string) []string {
	return samurai.Split(token, s.context, lists.Prefixes, lists.Suffixes)
}

// NewSamurai creates a new Samurai splitter that will work under the provided context.
func NewSamurai(local, global *samurai.FrequencyTable) step.Splitter {
	return samuraiSplitter{
		splitter: splitter{"samurai"},
		context:  samurai.NewTokenContext(local, global),
	}
}
