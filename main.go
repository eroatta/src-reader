package main

import (
	"fmt"
	"log"

	"github.com/eroatta/token/samurai"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/step"
)

func main() {
	newGoodMain("https://github.com/src-d/go-siva")
}

func newGoodMain(url string) {
	// stage: clone (and retrieve files)
	filesc, err := step.Clone(url)
	if err != nil {
		log.Fatalf("Error reading repository %s: %v", url, err)
	}

	parsedFiles := make([]code.File, 0)
	parsedc := step.Parse(filesc)
	for p := range parsedc {
		parsedFiles = append(parsedFiles, p)
	}

	// TODO: merge ASTs considering packages

	// for each package
	// apply the set of miners (preprocessors)
	frequencyMiner := "samurai"
	miningResults := step.Mine(parsedFiles, frequencyMiner)

	frequencyTable := samurai.NewFrequencyTable()
	freqc := miningResults[frequencyMiner]
	for input := range freqc {
		freq := input.(map[string]int)
		for token, count := range freq {
			//TODO: review
			if len(token) == 1 {
				continue
			}
			frequencyTable.SetOccurrences(token, count)
		}
	}
	log.Println(frequencyTable)

	// for each package (AST)
	// apply the set of splitters + expanders
	//numberProcessors := 1
	identc := step.Extract(parsedFiles)

	tCtx := samurai.NewTokenContext(frequencyTable, frequencyTable)
	splittedc := step.Split(identc, tCtx)
	expandedc := step.Expand(splittedc)

	for ident := range expandedc {
		log.Println("Identifier received")
		for alg, splits := range ident.Splits {
			log.Println(fmt.Sprintf("FuncDecl \"%s\" Splitted into: %v by %s", ident.Name, splits, alg))
		}

		for alg, expans := range ident.Expansions {
			log.Println(fmt.Sprintf("FuncDecl \"%s\" Expanded into: %v by %s", ident.Name, expans, alg))
		}
	}
}
