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

// Expand receives a code.Identifier and processes the available splits that
// can be expanded with the current algorithm.
// On AMAP, we rely on the related scoped declaration information for the identifier.
// If no decalaration information can be found, we avoid trying to expand the identifier
// because results can be broad.
func (a amapExpander) Expand(ident code.Identifier) []string {
	split, ok := ident.Splits[a.ApplicableOn()]
	if !ok {
		return []string{}
	}

	// TODO: use key
	scopedDecl, ok := a.scopedDeclarations[ident.Name]
	if !ok {
		return split
	}

	// TODO change strings.Join
	// TODO: also, we can use amap.NewTokenScope(code.ScopedDecl)
	scope := amap.NewTokenScope(scopedDecl.VariableDecls, scopedDecl.Name,
		strings.Join(scopedDecl.BodyText, " "), scopedDecl.Comments, scopedDecl.PackageComments)

	var expanded []string
	for _, token := range split {
		expanded = append(expanded, amap.Expand(token, scope, a.referenceText)...)
	}

	return expanded
}

func (a amapExpander) ApplicableOn() string {
	return "samurai"
}

// NewAMAP creates a new AMAP expander. It depends on scoped declarations and also on a
// reference text.
func NewAMAP(declarations map[string]miner.ScopedDecl, referenceText []string) step.Expander {
	return amapExpander{
		expander:           expander{"amap"},
		scopedDeclarations: declarations,
		referenceText:      referenceText,
	}
}
