package repositories

import (
	"log"
	"os"

	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Cloner provides the required functions to clone, get the list of files and folders that are part of, and
// retrieve a file from a repository.
type Cloner interface {
	Clone(url string) (*Repository, error)
	GetFiles(repo string) ([]os.FileInfo, error)
}

// DefaultGithubCloner builds a component that can access a GitHub repository, using the default configuration and client.
func DefaultGithubCloner() Cloner {
	return GoGitCloner{}
}

type GoGitCloner struct {
}

func (c GoGitCloner) Clone(url string) (*Repository, error) {
	return nil, nil
}

func (c GoGitCloner) GetFiles(repo string) ([]os.FileInfo, error) {
	return make([]os.FileInfo, 0), nil
}

// Clone clones a GitHub repository, and return a struct that represents it.
func Clone(url string) (*Repository, error) {
	log.Println("Cloning repository...")
	repository, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: url,
	})

	if err != nil {
		return nil, err
	}

	return &Repository{repository: repository}, nil
}

// Repository represents a remote GitHub repository.
type Repository struct {
	repository *git.Repository

	ID           string
	Hash         string
	Language     string
	Contributors int
}
