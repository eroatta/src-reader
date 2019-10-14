package cloner

import (
	"log"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/step"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

const (
	rootDir = ""
)

// New creates and initializes a new cloner.
func New() step.Cloner {
	return &goGitCloner{
		clonerFunc: goGitClonerFunc,
	}
}

type goGitCloner struct {
	clonerFunc clonerFunc
	repository *git.Repository
}

// goGitClonerFunc defines the interface for cloning a remote Git repository.
type clonerFunc func(url string) (*git.Repository, error)

// GoGitClonerFunc clones a remote GitHub repository using the src{d}/go-git client.
func goGitClonerFunc(url string) (*git.Repository, error) {
	return git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: url,
	})
}

func (c *goGitCloner) Clone(url string) (code.Repository, error) {
	log.Println("Cloning repository...")
	repository, err := c.clonerFunc(url)
	if err != nil {
		return code.Repository{}, err
	}
	c.repository = repository

	return code.Repository{Name: url}, nil
}

// Filenames retrieves the list of file names existing on a repository.
func (c *goGitCloner) Filenames() ([]string, error) {
	wt, err := c.repository.Worktree()
	if err != nil {
		return nil, err
	}

	return read(wt.Filesystem, rootDir)
}

func read(fs billy.Filesystem, rootDir string) ([]string, error) {
	files, err := fs.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			subDirFilenames, err := read(fs, fs.Join(rootDir, file.Name()))
			if err != nil {
				return nil, err
			}

			names = append(names, subDirFilenames...)
		} else {
			names = append(names, fs.Join(rootDir, file.Name()))
		}
	}

	return names, nil
}

// File provides the bytes representation of a given file.
func (c *goGitCloner) File(name string) ([]byte, error) {
	wt, err := c.repository.Worktree()
	if err != nil {
		return nil, err
	}

	fileInfo, err := wt.Filesystem.Stat(name)
	if err != nil {
		return nil, err
	}

	file, err := wt.Filesystem.Open(name)
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, fileInfo.Size())
	_, err = file.Read(bytes)

	return bytes, err
}
