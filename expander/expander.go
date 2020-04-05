package expander

import (
	"errors"

	"github.com/eroatta/src-reader/entity"
	log "github.com/sirupsen/logrus"
)

func NewExpanderFactory() entity.ExpanderAbstractFactory {
	return &expanderFactory{
		factories: map[string]entity.ExpanderFactory{
			"noexp": NewNoExpansionFactory(),
			"basic": NewBasicFactory(),
			"amap":  NewAMAPFactory(),
		},
	}
}

type expanderFactory struct {
	factories map[string]entity.ExpanderFactory
}

func (f expanderFactory) Get(name string) (entity.ExpanderFactory, error) {
	factory, ok := f.factories[name]
	if !ok {
		log.WithField("name", name).Error("no factory declared for the given name")
		return nil, errors.New("no factory defined")
	}

	return factory, nil
}

type expander struct {
	name string
}

// Name returns the name of the Expander.
func (e expander) Name() string {
	return e.name
}
