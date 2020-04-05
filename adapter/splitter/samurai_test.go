package splitter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eroatta/token/samurai"

	"github.com/eroatta/src-reader/adapter/splitter"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
)

func TestNewSamuraiFactory_ShouldReturnSamuraiSplitterFactory(t *testing.T) {
	factory := splitter.NewSamuraiFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnSamuraiFactory_WhenMissingLocalFrequencyTable_ShouldReturnError(t *testing.T) {
	staticInputs := map[string]interface{}{
		"GlobalFrequencyTable": samurai.NewFrequencyTable(),
	}
	miningResults := map[entity.MinerType]entity.Miner{}

	factory := splitter.NewSamuraiFactory()
	splitter, err := factory.Make(staticInputs, miningResults)

	assert.Nil(t, splitter)
	assert.Error(t, err)
}

func TestMake_OnSamuariFactory_WhenMissingGlobalFrequencyTable_ShouldReturnError(t *testing.T) {
	staticInputs := map[string]interface{}{}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.WordCount: miner.NewWordCount(),
	}

	factory := splitter.NewSamuraiFactory()
	splitter, err := factory.Make(staticInputs, miningResults)

	assert.Nil(t, splitter)
	assert.Error(t, err)
}

func TestSplit_OnSamurai_ShouldReturnAnArrayOfStrings(t *testing.T) {
	staticInputs := map[string]interface{}{
		"GlobalFrequencyTable": samurai.NewFrequencyTable(),
	}
	miningResults := map[entity.MinerType]entity.Miner{
		entity.WordCount: miner.NewWordCount(),
	}

	factory := splitter.NewSamuraiFactory()
	splitter, _ := factory.Make(staticInputs, miningResults)
	got := splitter.Split("car")

	assert.Equal(t, "samurai", splitter.Name())
	assert.ElementsMatch(t, []string{"car"}, got)
}
