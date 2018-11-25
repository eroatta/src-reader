package repositories

import (
	"errors"
	"log"
	"os"

	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// ClonerFunc defines the interface for cloning a remote Git repository.
type ClonerFunc func(url string) (*git.Repository, error)

// Clone takes a repository URL (https://github.com/USER/REPO or git@github.com:/USER/REPO) and
// clones it, using the provided cloner.
func Clone(cloner ClonerFunc, url string) (*git.Repository, error) {
	log.Println("Cloning repository...")
	repository, err := cloner(url)
	if err != nil {
		return nil, errors.New("Error cloning remote repository")
	}

	return repository, nil
}

// GoGitClonerFunc clones a remote GitHub repository using the src{d}/go-git client.
func GoGitClonerFunc(url string) (*git.Repository, error) {
	return git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: url,
	})
}

// FilesInfo retrieves the list of files and its related information on a given repository.
func FilesInfo(repository *git.Repository) ([]os.FileInfo, error) {
	return make([]os.FileInfo, 0), nil
}

// File retrieves the raw file as an array of bytes.
func File(repository *git.Repository, filename string) ([]byte, error) {
	return make([]byte, 0), nil
}
