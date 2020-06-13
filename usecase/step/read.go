package step

import (
	"context"
	"strings"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
)

// Read filters only .go files and reads them.
func Read(ctx context.Context, sc repository.SourceCodeRepository, location string, filenames []string) <-chan entity.File {
	namesc := make(chan string)
	go func() {
		for _, f := range filenames {
			if strings.HasSuffix(f, "_test.go") || !strings.HasSuffix(f, ".go") {
				continue
			}

			namesc <- f
		}

		close(namesc)
	}()

	filesc := make(chan entity.File)
	go func() {
		for n := range namesc {
			rawFile, err := sc.Read(ctx, location, n)

			file := entity.File{
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
