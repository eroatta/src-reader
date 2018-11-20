package repositories

import (
	"os"
)

// Cloner defines the required functions for cloning a repository, getting the list of files that belongs to the repository
// and retrieving those files.
type Cloner interface {
	Name() string
	Clone(url string) (*Repository, error)
	GetFilesInfo(repo string) ([]os.FileInfo, error)
	GetFile(filename string) ([]byte, error)
}

// DefaultGithubCloner provides an instance of the default implementation of the Cloner interface, capable
// of accessing a GitHub repository, using the default configuration and client.
func DefaultGithubCloner() Cloner {
	return GoGitCloner{}
}

type GoGitCloner struct {
}

//Name provides the name of the cloner implementation.
func (c GoGitCloner) Name() string {
	return "go-git"
}

// Clone takes a repository URL (https://github.com/USER/REPO or git@github.com:/USER/REPO) and
// clones it.
func (c GoGitCloner) Clone(url string) (*Repository, error) {
	return nil, nil
}

func (c GoGitCloner) GetFilesInfo(repo string) ([]os.FileInfo, error) {
	return make([]os.FileInfo, 0), nil
}

func (c GoGitCloner) GetFile(filename string) ([]byte, error) {
	return make([]byte, 0), nil
}
