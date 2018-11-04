package main

import (
	"log"
	"os"

	"gopkg.in/src-d/go-billy.v4/memfs"

	"github.com/eroatta/src-reader/url"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func main() {
	log.Println("Starting src-reader")

	uri := "https://github.com/src-d/go-siva"
	if !url.IsValidGithubRepoURL(uri) {
		log.Fatal("Invalid Repo URI")
	}

	repository, err := read(uri)
	if err != nil {
		log.Fatal("Error reading the repository")
	}

	log.Println("Repo read...")

	// get files from filesystem
	files, err := getFilesFromFileSystem(repository)
	if err != nil {
		log.Fatal("Error retrieving files...")
	}

	for _, file := range files {
		log.Println(file.Name())

	}
}

func read(url string) (*git.Repository, error) {
	repository, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: url,
	})

	return repository, err
}

func getFilesFromFileSystem(repository *git.Repository) ([]os.FileInfo, error) {
	wt, err := repository.Worktree()
	if err != nil {
		return nil, err
	}

	files, err := wt.Filesystem.ReadDir("")
	if err != nil {
		return nil, err
	}

	return files, nil
}
