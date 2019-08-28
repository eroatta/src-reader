package step_test

import (
	"go/ast"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/step"
	"github.com/stretchr/testify/assert"
)

func TestMine_OnNoFiles_ShouldReturnMinersWithoutResults(t *testing.T) {
	processed := step.Mine([]code.File{}, miner{name: "empty"})

	assert.Equal(t, 1, len(processed))

	emptyMiner, ok := processed["empty"].(miner)
	assert.True(t, ok)
	assert.NotNil(t, emptyMiner)
	assert.Equal(t, 0, emptyMiner.visits)
}

func TestMine_OnEmptyMiners_ShouldReturnNoResults(t *testing.T) {
	processed := step.Mine([]code.File{}, []step.Miner{}...)

	assert.Equal(t, 0, len(processed))
}

func TestMine_OnFileWithNilAST_ShouldReturnMinersWithoutResults(t *testing.T) {
	processed := step.Mine([]code.File{{Name: "main.go"}}, miner{name: "empty"})

	assert.Equal(t, 1, len(processed))

	emptyMiner, ok := processed["empty"].(miner)
	assert.True(t, ok)
	assert.NotNil(t, emptyMiner)
	assert.Equal(t, 0, emptyMiner.visits)
}

func TestMine_OnTwoMiners_ShouldReturnResultsForEveryMiner(t *testing.T) {
	assert.Fail(t, "not yet implemented")
}

type miner struct {
	name   string
	visits int
}

func (m miner) Name() string {
	return m.name
}

func (m miner) Visit(n ast.Node) ast.Visitor {
	m.visits++
	return m
}
