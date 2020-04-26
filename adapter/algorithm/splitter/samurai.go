package splitter

import (
	"errors"
	"strings"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/lists"
	"github.com/eroatta/token/samurai"
	log "github.com/sirupsen/logrus"
)

// NewSamuraiFactory creates a new Greedy splitter factory.
func NewSamuraiFactory() entity.SplitterFactory {
	return samuraiFactory{}
}

type samuraiFactory struct{}

func (f samuraiFactory) Make(miningResults map[string]entity.Miner) (entity.Splitter, error) {
	// build local frequency table from word count
	wordsMiner, ok := miningResults["wordcount"]
	if !ok {
		return nil, errors.New("unable to retrieve input from wordcount miner")
	}

	local := samurai.NewFrequencyTable()
	frequencies := wordsMiner.Results().(map[string]int)
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
	globalFreqTableMiner, ok := miningResults["global-frequency-table"]
	if !ok {
		return nil, errors.New("unable to retrieve input from global-frequency-table miner")
	}
	global := globalFreqTableMiner.Results().(*samurai.FrequencyTable)

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
func (s samuraiSplitter) Split(token string) []entity.Split {
	splits := []entity.Split{}
	for i, split := range strings.Split(samurai.Split(token, s.context, lists.Prefixes, lists.Suffixes), " ") {
		splits = append(splits, entity.Split{Order: i + 1, Value: split})
	}

	return splits
}
