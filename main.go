package main

import (
	"go/ast"
	"go/token"
	"log"

	"github.com/eroatta/token/greedy"
	"github.com/eroatta/token/lists"
	"github.com/eroatta/token/samurai"

	"github.com/eroatta/token/conserv"

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

	/*occurc, _ := mine(trees)
	freqTable := samurai.NewFrequencyTable()
	for input := range occurc {
		for token, count := range input {
			//TODO: review
			if len(token) == 1 {
				continue
			}
			freqTable.SetOccurrences(token, count)
		}
	}
	fmt.Println(freqTable)*/

	// for each package (AST)
	// apply the set of splitters + expanders
	//numberProcessors := 1
	/*identc := make(chan identifier)
	go func() {
		for _, t := range trees {
			ast.Inspect(t, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if ok {
					ident := identifier{
						file:       "main.go",
						position:   fn.Pos(),
						name:       fn.Name.String(),
						typ:        "FuncDecl",
						splits:     make(map[string][]string),
						expansions: make(map[string][]string),
					}
					identc <- ident
				}

				return true
			})
		}
		close(identc)
	}()

	log.Println("Closed identc...")

	tCtx := samurai.NewTokenContext(freqTable, freqTable)
	splittedc := split(identc, tCtx)
	expandedc := expand(splittedc)

	for ident := range expandedc {
		log.Println("Identifier received")
		for alg, splits := range ident.splits {
			log.Println(fmt.Sprintf("FuncDecl \"%s\" Splitted into: %v by %s", ident.name, splits, alg))
		}

		for alg, expans := range ident.expansions {
			log.Println(fmt.Sprintf("FuncDecl \"%s\" Expanded into: %v by %s", ident.name, expans, alg))
		}
	}*/
}

type identifier struct {
	file       string
	position   token.Pos
	name       string
	typ        string
	node       *ast.Node
	splits     map[string][]string
	expansions map[string][]string
}

func split(identc <-chan identifier, tCtx samurai.TokenContext) chan identifier {
	/*
		SPLITTING STEP: it splits all the identifiers in a AST, applying the given set of Splitters.
		Input: a list of *ast.File nodes, and a list of Splitters.
		Output: ?
	*/
	splittedc := make(chan identifier)
	go func() {
		for ident := range identc {
			log.Println("Splitting...")
			ident.splits["conserv"] = conserv.Split(ident.name)
			ident.splits["greedy"] = greedy.Split(ident.name, greedy.DefaultList)
			ident.splits["samurai"] = samurai.Split(ident.name, tCtx, lists.Prefixes, lists.Suffixes)

			splittedc <- ident
		}
		close(splittedc)
	}()

	return splittedc
}

func expand(identc <-chan identifier) chan identifier {
	expandedc := make(chan identifier)
	go func() {
		for ident := range identc {
			log.Println("Expanding...")
			// TODO: add expansions

			expandedc <- ident
		}

		close(expandedc)
	}()

	return expandedc
}
