package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/eroatta/src-reader/extractors"

	"github.com/eroatta/src-reader/repositories"

	"github.com/eroatta/src-reader/url"
)

func main() {
	log.Println("Starting src-reader")

	uri := "https://github.com/src-d/go-siva"
	if !url.IsValidGithubRepoURL(uri) {
		log.Fatal("Invalid Repo URI")
	}

	repository, err := repositories.Clone(repositories.GoGitClonerFunc, uri)
	if err != nil {
		log.Fatal("Error reading the repository")
	}

	log.Println("Repo read...")

	filenames, err := repositories.Filenames(repository)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Files: %v", filenames))

	samurai := extractors.NewSamuraiExtractor()
	fset := token.NewFileSet() // positions are relative to fset
	for _, name := range filenames {
		if !strings.HasSuffix(name, ".go") {
			continue
		}

		log.Println(fmt.Sprintf("Processing file: %s", name))

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

		//ast.Print(fset, node)
		extractors.Process(samurai, node)

		log.Println(fmt.Sprintf("Elements: %d", len(samurai.FreqTable())))
	}

	log.Println(fmt.Sprintf("Elements: %d", len(samurai.FreqTable())))
	for k, v := range samurai.FreqTable() {
		log.Println(fmt.Sprintf("Token: %s - Occurrencies: %d", k, v))
	}
}
