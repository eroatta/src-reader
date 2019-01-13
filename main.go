package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/eroatta/token-splitex/splitters"

	"github.com/eroatta/src-reader/extractors"

	"github.com/eroatta/src-reader/repositories"

	"github.com/eroatta/src-reader/url"
)

func main() {
	log.Println("Starting src-reader...")

	/*
		CLONING STEP: retrieves the source code from Github.
		It validates the given Github URI, clones the repository, checks for errors
		and also extract the list of file names that belong to the project.
		Input: Github URI.
		Output: a list of files and a filesystem so we can read them.
	*/
	log.Println("Beginning Cloning Step...")
	uri := "https://github.com/src-d/go-siva"
	if !url.IsValidGithubRepoURL(uri) {
		log.Fatal("Invalid Github repository URI.")
	}

	repository, err := repositories.Clone(repositories.GoGitClonerFunc, uri)
	if err != nil {
		log.Fatalf("Error reading repository %s: %v", uri, err)
	}

	filenames, err := repositories.Filenames(repository)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Repository %s includes the following files: %v", uri, filenames))

	/*
		PARSING STEP: parses a file and creates and AST so it can be explored lately.
		It filters non-go files and skips them, so we can focus only on Go source code.
		It also retrieves the binary files and parses it, using the ast/parser from Go.
		Input: a list of files and a filesystem so we can read them.
		Output: a list of ASTs, where each AST represents a file.
		Improvement opportunity: we could merge the ASTs into one collection of ast.Package.
	*/
	log.Println("Beginning Parsing Step...")

	var asts []*ast.File
	fset := token.NewFileSet() // positions are relative to fset
	for _, name := range filenames {
		if !strings.HasSuffix(name, ".go") {
			continue
		}

		log.Println(fmt.Sprintf("Parsing file: %s", name))

		rawFile, err := repositories.File(repository, name)
		if err != nil {
			log.Fatal(err)
			continue
		}

		node, err := parser.ParseFile(fset, name, rawFile, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
			continue
		}

		asts = append(asts, node)
	}

	/*
		MINING STEP: extracts all the information required by each splitting and expanding algorithm.
		It traverses each AST with every defined miner, so they can extract the required info.
		Input: a list of *ast.File nodes, and a list of Miners (new interface?)
		Output: ?
	*/
	log.Println("Beginning Mining Step...")
	samurai := extractors.NewSamuraiExtractor()
	for _, ast := range asts {
		extractors.Process(samurai, ast)

		log.Println(fmt.Sprintf("Elements after processing a new AST: %d", len(samurai.FreqTable())))
	}

	freqTable := buildFrequencyTable(samurai.FreqTable())
	log.Println(fmt.Sprintf("Program Frequency Table - Total Occurencies: %d", freqTable.TotalOccurrences()))
	log.Println(fmt.Sprintf("Frequency for %s: %f", "index", freqTable.Frequency("index")))

	/*
		SPLITTING STEP: it splits all the identifiers in a AST, applying the given set of Splitters.
		Input: a list of *ast.File nodes, and a list of Splitters.
		Output: ?
	*/
	samuraiSplitter := splitters.NewSamurai(freqTable, freqTable, nil, nil)
	splits, err := samuraiSplitter.Split("srccode")
	if err != nil {
		log.Fatalf("Unable to split token \"%s\": %v", "srccode", err)
	}

	log.Println(fmt.Sprintf("Splits for token \"%s\": %v", "srccode", splits))
}

func buildFrequencyTable(input map[string]int) *splitters.FrequencyTable {
	freqTable := splitters.NewFrequencyTable()
	for token, count := range input {
		freqTable.SetOccurrences(token, count)
	}

	return freqTable
}
