package repositories

import (
	"os"
)

// Cloner defines the required functions for cloning a repository, getting the list of files that belongs to the repository
// and retrieving those files.
type Cloner interface {
	Name() string
	Clone(url string) (*Repository, error)
	GetFilesInfo() ([]os.FileInfo, error)
	GetFile(filename string) ([]byte, error)
}

// DefaultGithubCloner provides an instance of the default implementation of the Cloner interface, capable
// of accessing a GitHub repository, using the default configuration and client.
func DefaultGithubCloner() Cloner {
	return GoGitCloner{}
}

// GoGitCloner stores information and provides the required methods to implement the Cloner interface
// for the src{d}/go-git GitHub client.
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

// GetFilesInfo retrieves the list of files on a given repository.
func (c GoGitCloner) GetFilesInfo() ([]os.FileInfo, error) {
	return make([]os.FileInfo, 0), nil
}

// GetFile retrieves the raw file as an array of bytes.
func (c GoGitCloner) GetFile(filename string) ([]byte, error) {
	return make([]byte, 0), nil
}
