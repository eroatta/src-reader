package step

import (
	"context"
	"strings"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/repository"
)

func Read(ctx context.Context, sc repository.SourceCodeRepository, location string, filenames []string) <-chan code.File {
	namesc := make(chan string)
	go func() {
		for _, f := range filenames {
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
			rawFile, err := sc.Read(ctx, location, n)

			file := code.File{
				Name:  n,
				Raw:   rawFile,
				Error: err,
			}
			filesc <- file
		}

		close(filesc)
	}()

	return filesc
}
