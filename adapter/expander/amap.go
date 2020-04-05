package expander

import (
	"fmt"
	"strings"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/miner"
	"github.com/eroatta/token/amap"
)

// NewAMAPFactory creates a new AMAP expanders factory.
func NewAMAPFactory() entity.ExpanderFactory {
	return amapFactory{}
}

type amapFactory struct{}

func (f amapFactory) Make(miningResults map[entity.MinerType]entity.Miner) (entity.Expander, error) {
	declarationsMiner, ok := miningResults[entity.MinerScopedDeclarations]
	if !ok {
		return nil, fmt.Errorf("unable to retrieve input from %s", entity.MinerScopedDeclarations)
	}

	scopedDeclarations := declarationsMiner.(miner.Scope).ScopedDeclarations()
	referenceText := []string{} // TODO: add reference text

	return &amapExpander{
		expander:           expander{"amap"},
		scopedDeclarations: scopedDeclarations,
		referenceText:      referenceText,
	}, nil
}

type amapExpander struct {
	expander
	scopedDeclarations map[string]entity.ScopedDecl
	referenceText      []string
}

// Expand receives a entity.Identifier and processes the available splits that
// can be expanded with the current algorithm.
// On AMAP, we rely on the related scoped declaration information for the identifier.
// If no decalaration information can be found, we avoid trying to expand the identifier
// because results can be broad.
func (a amapExpander) Expand(ident entity.Identifier) []string {
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
