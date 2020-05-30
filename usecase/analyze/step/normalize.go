package step

import (
	"github.com/eroatta/src-reader/entity"
)

// Normalize returns a channel of entity.Identifier where each element has been normalized.
func Normalize(identc chan entity.Identifier) chan entity.Identifier {
	normalizedc := make(chan entity.Identifier)
	go func() {
		for ident := range identc {
			ident.Normalize()
			normalizedc <- ident
		}

		close(normalizedc)
	}()

	return normalizedc
}
