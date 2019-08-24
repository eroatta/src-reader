package step

import (
	"go/ast"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/miner"
)

// Miner interface is used to define a custom miner.
type Miner interface {
	// Name provides the name of the miner.
	Name() string
	// Visit applies the mining logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
	// Results returns the results once the mining is completed.
	Results() interface{}
}

// Mine traverses each Abstract Syntax Tree to apply every given miner to extract
// the required pre-processing information.
func Mine(parsed []code.File, miners ...string) map[string]chan interface{} {
	results := make(map[string]chan interface{})
	for _, minerName := range miners {
		resc := make(chan interface{}, 0)
		go func() {
			for _, f := range parsed {
				miner := miner.New(minerName).(Miner)
				ast.Walk(miner, f.AST)
				resc <- miner.Results()
			}

			close(resc)
		}()

		results[minerName] = resc
	}

	return results
}
