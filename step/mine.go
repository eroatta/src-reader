package step

import (
	"go/ast"
	"sync"

	"github.com/eroatta/src-reader/code"
)

type MinerType string

type MiningResults map[MinerType]interface{}

const (
	Words                MinerType = "Words"
	Phrases              MinerType = "Phrases"
	LocalFrequencyTable  MinerType = "Local Frequency Table"
	GlobalFrequencyTable MinerType = "Global Frequency Table"
	ScopedDeclarations   MinerType = "Scoped Declarations"
)

// Miner interface is used to define a custom miner.
type Miner interface {
	// Name provides the name of the miner.
	Name() string
	// Visit applies the mining logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
}

// Mine traverses each Abstract Syntax Tree and applies every given miner to extract
// the required pre-processing information. It returns a map of miners after work is done.
func Mine(parsed []code.File, miners ...Miner) map[string]Miner {
	minersc := make(chan Miner)

	var wg sync.WaitGroup
	wg.Add(len(miners))
	for _, miner := range miners {
		go func(miner Miner) {
			defer wg.Done()

			for _, f := range parsed {
				if f.AST == nil {
					continue
				}

				ast.Walk(miner, f.AST)
			}

			minersc <- miner
		}(miner)
	}

	go func() {
		wg.Wait()
		close(minersc)
	}()

	results := make(map[string]Miner)
	for miner := range minersc {
		results[miner.Name()] = miner
	}

	return results
}
