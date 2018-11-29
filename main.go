package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"

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
	for _, name := range filenames {
		log.Println(name)
	}

	rawFile, err := repositories.File(repository, "common.go")
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "common.go", rawFile, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}
