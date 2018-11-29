package main

import (
	"fmt"
	"log"

	"github.com/eroatta/src-reader/repositories"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

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

	head, err := repository.Head()
	commit, err := repository.CommitObject(head.Hash())
	tree, err := commit.Tree()
	tree.Files().ForEach(func(f *object.File) error {
		fmt.Printf("100644 blob %s    %s\n", f.Hash, f.Name)
		return nil
	})
	// get files from filesystem
	files, err := repositories.FilesInfo(repository)
	if err != nil {
		log.Fatal("Error retrieving files...")
	}

	for _, file := range files {
		log.Println(file.Name())
	}
}
