package expander

import (
	"strings"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/src-reader/step"
	"github.com/eroatta/token/amap"
)

type amapExpander struct {
	expander
	scopedDeclarations map[string]miner.ScopedDecl
	referenceText      []string
}

func (a amapExpander) Expand(ident code.Identifier) []string {
	scopedDecl, ok := a.scopedDeclarations[ident.Name]
	if !ok {
		// TODO perhaps we should return the identifier split
		return []string{}
	}

	// TODO change strings.Join
	scope := amap.NewTokenScope(scopedDecl.VariableDecls, scopedDecl.Name,
		strings.Join(scopedDecl.BodyText, " "), scopedDecl.Comments, scopedDecl.PackageComments)

	var expanded []string
	tokens := ident.Splits[a.ApplicableOn()]
	for _, token := range tokens {
		expansions := amap.Expand(token, scope, a.referenceText)
		expanded = append(expanded, expansions...)
	}

	return expanded
}

func (a amapExpander) ApplicableOn() string {
	return "samurai"
}

// TODO check if ScopedDecl should be part of miner...
func NewAMAP(declarations map[string]miner.ScopedDecl) step.Expander {
	return amapExpander{
		expander:           expander{"amap"},
		scopedDeclarations: declarations,
		referenceText:      []string{}, // TODO where do we get it?
	}
}
