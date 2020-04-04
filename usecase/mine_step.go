package usecase

import (
	"go/ast"
	"sync"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
)

type MiningResults map[entity.MinerType]interface{}

// mine traverses each Abstract Syntax Tree and applies every given miner to extract
// the required pre-processing information. It returns a map of miners after work is done.
func mine(parsed []code.File, miners ...entity.Miner) map[entity.MinerType]entity.Miner {
	minersc := make(chan entity.Miner)

	var wg sync.WaitGroup
	wg.Add(len(miners))
	for _, miner := range miners {
		go func(miner entity.Miner) {
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

	results := make(map[entity.MinerType]entity.Miner)
	for miner := range minersc {
		results[miner.Type()] = miner
	}

	return results
}
