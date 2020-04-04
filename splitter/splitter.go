package splitter

import (
	"errors"

	"github.com/eroatta/src-reader/entity"
	log "github.com/sirupsen/logrus"
)

func NewSplitterFactory() entity.SplitterAbstractFactory {
	return &splitterFactory{
		factories: map[string]entity.SplitterFactory{
			"conserv": NewConservFactory(),
		},
	}
}

type splitterFactory struct {
	factories map[string]entity.SplitterFactory
}

func (f splitterFactory) Get(name string) (entity.SplitterFactory, error) {
	factory, ok := f.factories[name]
	if !ok {
		log.WithField("name", name).Error("no factory declared for the given name")
		return nil, errors.New("no factory defined")
	}

	return factory, nil
}

type splitter struct {
	name string
}

// Name returns the name of the splitter.
func (s splitter) Name() string {
	return s.name
}
