package repositories

import (
	"errors"
	"log"
	"os"

	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

const (
	rootDir = ""
)

var (
	// ErrCloningRepository represents the default error message when an error occurs while cloning a repository.
	ErrCloningRepository = errors.New("Error cloning the remote repository")
)

// ClonerFunc defines the interface for cloning a remote Git repository.
type ClonerFunc func(url string) (*git.Repository, error)

// Clone takes a repository URL (https://github.com/USER/REPO or git@github.com:/USER/REPO) and
// clones it, using the provided cloner.
func Clone(cloner ClonerFunc, url string) (*git.Repository, error) {
	log.Println("Cloning repository...")
	repository, err := cloner(url)
	if err != nil {
		return nil, ErrCloningRepository
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
	wt, err := repository.Worktree()
	if err != nil {
		return nil, err
	}

	return read(wt.Filesystem, rootDir)
}

func read(fs billy.Filesystem, rootDir string) ([]os.FileInfo, error) {
	files, err := fs.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	filesInfo := make([]os.FileInfo, 0)
	for _, file := range files {
		if file.IsDir() {
			subFiles, err := read(fs, file.Name())
			if err != nil {
				return nil, err
			}

			filesInfo = append(filesInfo, subFiles...)
		} else {
			filesInfo = append(filesInfo, file)
		}
	}

	return filesInfo, nil
}

// File retrieves the raw file as an array of bytes.
func File(repository *git.Repository, filename string) ([]byte, error) {
	return make([]byte, 0), nil
}
