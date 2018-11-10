package repositories

import (
	"log"

	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

//func Clone(cloner Cloner, url string) (*Repository, error) {
//
//}

type Cloner interface {
	Clone(url string) (*Repository, error)
}

func DefaultGithubCloner() Cloner {
	return GoGitCloner{}
}

type GoGitCloner struct {
}

func (c GoGitCloner) Clone(url string) (*Repository, error) {
	return nil, nil
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
}
