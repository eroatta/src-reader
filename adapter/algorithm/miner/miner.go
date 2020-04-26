package miner

import (
	"errors"

	"github.com/eroatta/src-reader/entity"
	log "github.com/sirupsen/logrus"
)

func NewMinerFactory() entity.MinerAbstractFactory {
	return &minerFactory{
		factories: map[string]entity.MinerFactory{
			"comments":               NewCommentsFactory(),
			"declarations":           NewDeclarationsFactory(),
			"global-frequency-table": NewGlobalFreqTableFactory(),
			"scoped-declarations":    NewScopesFactory(),
			"wordcount":              NewWordcountFactory(),
		},
	}
}

type minerFactory struct {
	factories map[string]entity.MinerFactory
}

// Get retrieves the miner factory for the requested algorithm.
func (f minerFactory) Get(name string) (entity.MinerFactory, error) {
	factory, ok := f.factories[name]
	if !ok {
		log.WithField("name", name).Error("no factory declared for the given name")
		return nil, errors.New("no factory defined")
	}

	return factory, nil
}

type miner struct {
	name string
}

// Name return the name of the miner.
func (m miner) Name() string {
	return m.name
}
