package main

import (
	"fmt"
	"log"

	"github.com/eroatta/token/conserv"
	"github.com/eroatta/token/greedy"
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
	// cloning step
	_, filesc, err := step.Clone(url, cloner.New())
	if err != nil {
		log.Fatalf("Error reading repository %s: %v", url, err)
	}

	// parsing step
	parsedc := step.Parse(filesc)
	files := step.Merge(parsedc)

	// mining step
	countMiner := miner.NewCount()
	declarationMiner := miner.NewDeclaration(lists.Dicctionary)
	scopeMiner := miner.NewScope("fail")

	miningResults := step.Mine(files, countMiner, declarationMiner, scopeMiner)

	frequencyTable := samurai.NewFrequencyTable()

	countResults := miningResults[countMiner.Name()].(miner.Count)
	freq := countResults.Results().(map[string]int)
	for token, count := range freq {
		if len(token) == 1 {
			continue
		}
		frequencyTable.SetOccurrences(token, count)
	}

	declarationResults := miningResults[declarationMiner.Name()].(miner.Declaration)
	decls := declarationResults.Decls()
	for k := range decls {
		log.Println(fmt.Sprintf("Declaration: %s", k))
	}

	scopeResults := miningResults[scopeMiner.Name()].(miner.Scope)
	scopes := scopeResults.ScopedDeclarations()
	for k := range scopes {
		log.Println(fmt.Sprintf("Scope: %s", k))
	}

	// extraction step
	identc := step.Extract(files, extractor.New)

	// splitting step
	tCtx := samurai.NewTokenContext(frequencyTable, frequencyTable)
	samuraiSplitter := newSamuraiSplitter(tCtx)
	conservSplitter := newConservSplitter()
	greedySplitter := newGreedySplitter()

	splittedc := step.Split(identc, samuraiSplitter, conservSplitter, greedySplitter)

	// expansion step
	expandedc := step.Expand(splittedc)

	// storing step
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

func newConservSplitter() splitter {
	return splitter{
		name: "conserv",
		fn: func(t string) []string {
			return conserv.Split(t)
		},
	}
}

func newGreedySplitter() splitter {
	return splitter{
		name: "greedy",
		fn: func(t string) []string {
			return greedy.Split(t, greedy.DefaultList)
		},
	}
}
