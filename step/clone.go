package step

import (
	"strings"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/repositories"
	"gopkg.in/src-d/go-git.v4"
)

// Clone retrieves the source code from GitHub, based on a given URL.
// It access the repository, clones it, filters non-go files and returns
// a channel of File elements.
func Clone(url string) (<-chan code.File, error) {
	repository, err := repositories.Clone(repositories.GoGitClonerFunc, url)
	if err != nil {
		return nil, err
	}

	files, err := repositories.Filenames(repository)
	if err != nil {
		return nil, err
	}

	out := make(chan string)
	go func() {
		for _, f := range files {
			if !strings.HasSuffix(f, ".go") {
				continue
			}
			out <- f
		}
		close(out)
	}()

	return retrieve(repository, out), nil
}

func retrieve(repo *git.Repository, namesc <-chan string) chan code.File {
	filesc := make(chan code.File)
	go func() {
		for n := range namesc {
			rawFile, err := repositories.File(repo, n)
			// TODO: review errors (do I need error channel?)
			if err != nil {
				continue
			}

			file := code.File{
				Name: n,
				Raw:  rawFile,
			}
			filesc <- file
		}
		close(filesc)
	}()

	return filesc
}
