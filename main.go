package main

import (
	"log"

	"github.com/eroatta/token/lists"
	"github.com/eroatta/token/samurai"

	"github.com/eroatta/src-reader/cloner"
	"github.com/eroatta/src-reader/extractor"
	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/src-reader/step"
	"github.com/eroatta/src-reader/storer"
)

func main() {
	newGoodMain("https://github.com/src-d/go-siva")
}

func newGoodMain(url string) {
	_, filesc, err := step.Clone(url, cloner.New())
	if err != nil {
		log.Fatalf("Error reading repository %s: %v", url, err)
	}

	parsedc := step.Parse(filesc)
	files := step.Merge(parsedc)

	countMiner := miner.NewCount()
	miningResults := step.Mine(files, countMiner)

	frequencyTable := samurai.NewFrequencyTable()

	countResults := miningResults[countMiner.Name()].(miner.Count)
	freq := countResults.Results().(map[string]int)
	for token, count := range freq {
		if len(token) == 1 {
			continue
		}
		frequencyTable.SetOccurrences(token, count)
	}

	identc := step.Extract(files, extractor.New)

	tCtx := samurai.NewTokenContext(frequencyTable, frequencyTable)
	samuraiSplitter := newSamuraiSplitter(tCtx)

	splittedc := step.Split(identc, samuraiSplitter)
	expandedc := step.Expand(splittedc)
	errors := step.Store(expandedc, storer.New())
	if len(errors) > 0 {
		log.Fatal("Something failed")
	}
}

type splitter struct {
	name string
	fn   func(string) []string
}

func (s splitter) Name() string {
	return s.name
}

func (s splitter) Split(token string) []string {
	return s.fn(token)
}

func newSamuraiSplitter(tokenContext samurai.TokenContext) splitter {
	return splitter{
		name: "samurai",
		fn: func(t string) []string {
			return samurai.Split(t, tokenContext, lists.Prefixes, lists.Suffixes)
		},
	}
}
