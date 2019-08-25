package main

import (
	"fmt"
	"log"

	"github.com/eroatta/token/lists"
	"github.com/eroatta/token/samurai"
	"gopkg.in/src-d/go-git.v4"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/step"
)

func main() {
	newGoodMain("https://github.com/src-d/go-siva")
}

func newGoodMain(url string) {
	_, filesc, err := step.Clone(url, &cloner{})
	if err != nil {
		log.Fatalf("Error reading repository %s: %v", url, err)
	}

	parsedc := step.Parse(filesc)
	files := step.Merge(parsedc)

	samuraiMiner := miner.NewSamuraiExtractor() // miner.New("samurai").(step.Miner)
	miningResults := step.Mine(files, samuraiMiner)

	frequencyTable := samurai.NewFrequencyTable()

	samuraiResults := miningResults[samuraiMiner.Name()].(miner.SamuraiExtractor)
	freq := samuraiResults.Results().(map[string]int)
	for token, count := range freq {
		//TODO: review
		if len(token) == 1 {
			continue
		}
		frequencyTable.SetOccurrences(token, count)
	}

	log.Println(frequencyTable)

	// for each package (AST)
	// apply the set of splitters + expanders

	identc := step.Extract(files)

	tCtx := samurai.NewTokenContext(frequencyTable, frequencyTable)
	samuraiSplitter := newSamuraiSplitter(tCtx)

	splittedc := step.Split(identc, samuraiSplitter)
	expandedc := step.Expand(splittedc)
	// sink or consumer for the expanded identifiers
	// results :=

	errors := step.Store(expandedc, mydb{})
	if len(errors) > 0 {
		log.Fatal("Something failed")
	}

	// update process information
}

type cloner struct {
	repo *git.Repository
}

func (c *cloner) Clone(url string) (code.Repository, error) {
	repo, err := repository.Clone(repository.GoGitClonerFunc, url)
	if err != nil {
		return code.Repository{}, err
	}
	c.repo = repo

	return code.Repository{Name: url}, nil
}

// Filenames retrieves the list of file names existing on a repository.
func (c *cloner) Filenames() ([]string, error) {
	return repository.Filenames(c.repo)
}

// File provides the bytes representation of a given file.
func (c *cloner) File(name string) ([]byte, error) {
	return repository.File(c.repo, name)
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

type mydb struct {
}

func (m mydb) Save(ident code.Identifier) error {
	log.Println("Storing identifier...")
	for alg, splits := range ident.Splits {
		log.Println(fmt.Sprintf("FuncDecl \"%s\" Splitted into: %v by %s", ident.Name, splits, alg))
	}

	for alg, expans := range ident.Expansions {
		log.Println(fmt.Sprintf("FuncDecl \"%s\" Expanded into: %v by %s", ident.Name, expans, alg))
	}

	return nil
}
