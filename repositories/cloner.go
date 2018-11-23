package repositories

import (
	"log"

	"gopkg.in/src-d/go-billy.v4"

	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Cloner defines the required functions for cloning a repository, getting the list of files that belongs to the repository
// and retrieving those files.
type Cloner interface {
	Name() string
	Clone(url string) (Repository, error)
}

// GoGitCloner stores information and provides the required methods to implement the Cloner interface
// for the src{d}/go-git GitHub client.
type GoGitCloner struct {
	storage    *memory.Storage
	filesystem billy.Filesystem
}

// DefaultGithubCloner provides an instance of the default implementation of the Cloner interface, capable
// of accessing a GitHub repository, using the default configuration and client.
func DefaultGithubCloner() Cloner {
	return GoGitCloner{
		storage:    memory.NewStorage(),
		filesystem: memfs.New(),
	}
}

// Name provides the name of the cloner implementation.
func (c GoGitCloner) Name() string {
	return "go-git"
}

// Clone takes a repository URL (https://github.com/USER/REPO or git@github.com:/USER/REPO) and
// clones it.
func (c GoGitCloner) Clone(url string) (Repository, error) {
	log.Println("Cloning repository...")
	repository, err := git.Clone(c.storage, c.filesystem, &git.CloneOptions{
		URL: url,
	})

	return GoGitRepository{repository: repository}, err
}
