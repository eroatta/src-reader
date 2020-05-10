package splitter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/port/outgoing/adapter/algorithm/miner"
	"github.com/eroatta/src-reader/port/outgoing/adapter/algorithm/splitter"
	"github.com/eroatta/token/samurai"
)

func TestNewSamuraiFactory_ShouldReturnSamuraiSplitterFactory(t *testing.T) {
	factory := splitter.NewSamuraiFactory()

	assert.NotNil(t, factory)
}

func TestMake_OnSamuraiFactory_WhenMissingLocalFrequencyTable_ShouldReturnError(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"global-frequency-table": miner.NewGlobalFreqTable(nil),
	}

	factory := splitter.NewSamuraiFactory()
	splitter, err := factory.Make(miningResults)

	assert.Nil(t, splitter)
	assert.Error(t, err)
}

func TestMake_OnSamuraiFactory_WhenMissingGlobalFrequencyTable_ShouldReturnError(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"wordcount": miner.NewWordCount(),
	}

	factory := splitter.NewSamuraiFactory()
	splitter, err := factory.Make(miningResults)

	assert.Nil(t, splitter)
	assert.Error(t, err)
}

func TestSplit_OnSamurai_ShouldReturnSplit(t *testing.T) {
	miningResults := map[string]entity.Miner{
		"wordcount":              miner.NewWordCount(),
		"global-frequency-table": miner.NewGlobalFreqTable(samurai.NewFrequencyTable()),
	}

	factory := splitter.NewSamuraiFactory()
	splitter, _ := factory.Make(miningResults)
	got := splitter.Split("car")

	assert.Equal(t, "samurai", splitter.Name())
	assert.Equal(t, []entity.Split{{Order: 1, Value: "car"}}, got)
}
