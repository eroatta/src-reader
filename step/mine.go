package step

import (
	"go/ast"
	"sync"

	"github.com/eroatta/src-reader/code"
)

// Miner interface is used to define a custom miner.
type Miner interface {
	// Name provides the name of the miner.
	Name() string
	// Visit applies the mining logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
}

// Mine traverses each Abstract Syntax Tree to apply every given miner to extract
// the required pre-processing information.
func Mine(parsed []code.File, miners ...Miner) map[string]Miner {
	minersc := make(chan Miner, 0)

	var wg sync.WaitGroup
	wg.Add(len(miners))
	for _, miner := range miners {
		go func(miner Miner) {
			defer wg.Done()

			for _, f := range parsed {
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
