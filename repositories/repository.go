package repositories

import (
	"os"

	git "gopkg.in/src-d/go-git.v4"
)

// Repository represents a remote repository.
type Repository interface {
	GetFilesInfo() ([]os.FileInfo, error)
	GetFile(filename string) ([]byte, error)
}

// GoGitRepository represents a remote GitHub repository.
type GoGitRepository struct {
	//Repository

	repository *git.Repository
}

// GetFilesInfo retrieves the list of files on a given repository.
func (r GoGitRepository) GetFilesInfo() ([]os.FileInfo, error) {
	return make([]os.FileInfo, 0), nil
}

// GetFile retrieves the raw file as an array of bytes.
func (r GoGitRepository) GetFile(filename string) ([]byte, error) {
	return make([]byte, 0), nil
}
