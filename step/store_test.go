package step_test

import (
	"errors"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/step"
	"github.com/stretchr/testify/assert"
)

func TestStore_OnClosedChannel_ShouldStoreNoElements(t *testing.T) {
	identc := make(chan code.Identifier)
	close(identc)

	storer := &storer{}
	failed := step.Store(identc, storer)

	assert.Empty(t, failed)
	assert.Equal(t, 0, storer.calls)
}

func TestStore_OnFailingStorer_ShouldReturnAllElementsWithErrors(t *testing.T) {
	identc := make(chan code.Identifier)
	go func() {
		identc <- code.Identifier{
			Name: "crtfile",
			Splits: map[string][]string{
				"test": []string{"crt", "file"},
			},
			Expansions: map[string][]string{
				"test": []string{"control", "delete"},
			},
		}
		close(identc)
	}()

	storer := &storer{err: errors.New("aborted connection")}
	failed := step.Store(identc, storer)

	assert.Equal(t, 1, len(failed))
	assert.EqualError(t, failed[0].Error, "aborted connection")

	assert.Equal(t, 1, storer.calls)
}

func TestStore_OnWorkingStorer_ShouldReturnNoElements(t *testing.T) {
	identc := make(chan code.Identifier)
	go func() {
		identc <- code.Identifier{
			Name: "crtfile",
			Splits: map[string][]string{
				"test": []string{"crt", "file"},
			},
			Expansions: map[string][]string{
				"test": []string{"control", "delete"},
			},
		}
		close(identc)
	}()

	storer := &storer{}
	failed := step.Store(identc, storer)

	assert.Equal(t, 0, len(failed))
	assert.Equal(t, 1, storer.calls)
}

type storer struct {
	err   error
	calls int
}

func (s *storer) Save(ident code.Identifier) error {
	s.calls++
	return s.err
}
