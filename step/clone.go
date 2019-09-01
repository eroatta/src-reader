package step

import (
	"strings"

	"github.com/eroatta/src-reader/code"
)

// Cloner interface is used to define a custom cloner.
type Cloner interface {
	// Clone accesses a repository and clones it.
	Clone(url string) (code.Repository, error)
	// Filenames retrieves the names of the existing files on a repository.
	Filenames() ([]string, error)
	// File provides the bytes representation of a given file.
	File(name string) ([]byte, error)
}

// Clone retrieves the source code from a given URL. It access the repository, clones it,
// filters non-go files and returns a channel of code.File elements.
func Clone(url string, cloner Cloner) (*code.Repository, <-chan code.File, error) {
	repo, err := cloner.Clone(url)
	if err != nil {
		return nil, nil, err
	}

	files, err := cloner.Filenames()
	if err != nil {
		return nil, nil, err
	}

	namesc := make(chan string)
	go func() {
		for _, f := range files {
			if !strings.HasSuffix(f, ".go") {
				continue
			}

			namesc <- f
		}

		close(namesc)
	}()

	filesc := make(chan code.File)
	go func() {
		for n := range namesc {
			rawFile, err := cloner.File(n)

			file := code.File{
				Name:  n,
				Raw:   rawFile,
				Error: err,
			}
			filesc <- file
		}

		close(filesc)
	}()

	return &repo, filesc, nil
}
