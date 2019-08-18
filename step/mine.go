package step

import (
	"go/ast"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/miner"
)

// Miner TODO
type Miner interface {
	Name() string
	Visit(node ast.Node) ast.Visitor
	Results() interface{}
}

// Mine TODO
func Mine(parsed []code.File, miners ...string) map[string]chan interface{} {
	/*
		MINING STEP: extracts all the information required by each splitting and expanding algorithm.
		It traverses each AST with every defined miner, so they can extract the required info.
		Input: a list of *ast.File nodes, and a list of Miners (new interface?)
		Output: ?
	*/
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

/*
occurc := make(chan map[string]int)
	go func() {
		for _, t := range trees {
			ext := extractors.NewSamuraiExtractor()
			ast.Walk(ext, t)
			occurc <- ext.FreqTable()
		}
		close(occurc)
	}()

	scopesc := make(chan map[string]amap.TokenScope)
	go func() {
		for _, t := range trees {
			ext := extractors.NewAmap("main.go")
			ast.Walk(ext, t)
			scopesc <- ext.Scopes()
		}
		close(scopesc)
	}()
*/
