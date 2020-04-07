package splitter

import (
	"fmt"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/token/lists"
	"github.com/eroatta/token/samurai"
	log "github.com/sirupsen/logrus"
)

// NewSamuraiFactory creates a new Greedy splitter factory.
func NewSamuraiFactory() entity.SplitterFactory {
	return samuraiFactory{}
}

type samuraiFactory struct{}

func (f samuraiFactory) Make(miningResults map[entity.MinerType]entity.Miner) (entity.Splitter, error) {
	// build local frequency table from word count
	wordsMiner, ok := miningResults[entity.MinerWordCount]
	if !ok {
		return nil, fmt.Errorf("unable to retrieve input from %s", entity.MinerWordCount)
	}

	local := samurai.NewFrequencyTable()
	frequencies := wordsMiner.(miner.WordCount).Results()
	for token, count := range frequencies {
		if len(token) == 1 {
			continue
		}

		err := local.SetOccurrences(token, count)
		if err != nil {
			log.WithField(token, count).Warn("unable to include token on local frequency table")
			continue
		}
	}

	// extract global frequency table
	globalFreqTableMiner, ok := miningResults[entity.MinerGlobalFrequencyTable]
	if !ok {
		return nil, fmt.Errorf("unable to retrieve input from %s", entity.MinerGlobalFrequencyTable)
	}
	global := globalFreqTableMiner.(miner.GlobalFreqTable).Table()

	return samuraiSplitter{
		splitter: splitter{"samurai"},
		context:  samurai.NewTokenContext(local, global),
	}, nil
}

type samuraiSplitter struct {
	splitter
	context samurai.TokenContext
}

// Split splits a token using the Samurai splitter.
func (s samuraiSplitter) Split(token string) []string {
	return samurai.Split(token, s.context, lists.Prefixes, lists.Suffixes)
}
