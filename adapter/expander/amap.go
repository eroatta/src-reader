package expander

import (
	"errors"
	"strings"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/token/amap"
)

// NewAMAPFactory creates a new AMAP expanders factory.
func NewAMAPFactory() entity.ExpanderFactory {
	return amapFactory{}
}

type amapFactory struct{}

func (f amapFactory) Make(miningResults map[string]entity.Miner) (entity.Expander, error) {
	declarationsMiner, ok := miningResults["scoped-declarations"]
	if !ok {
		return nil, errors.New("unable to retrieve input from scoped-declartions miner")
	}
	scopedDeclarations := declarationsMiner.Results().(map[string]entity.ScopedDecl)

	commentsMiner, ok := miningResults["comments"]
	if !ok {
		return nil, errors.New("unable to retrieve input from comments miner")
	}
	referenceText := commentsMiner.Results().([]string)

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
func (a amapExpander) Expand(ident entity.Identifier) []entity.Expansion {
	splits, ok := ident.Splits[a.ApplicableOn()]
	if !ok {
		return []entity.Expansion{}
	}

	declarationID := ident.ID
	if ident.IsLocal() {
		declarationID = ident.Parent
	}
	scopedDecl, ok := a.scopedDeclarations[declarationID]
	if !ok {
		expansions := make([]entity.Expansion, len(splits))
		for i, split := range splits {
			expansions[i] = entity.Expansion{From: split.Value, Values: []string{split.Value}}
		}
		return expansions
	}

	// TODO change strings.Join
	scope := amap.NewTokenScope(scopedDecl.VariableDecls, scopedDecl.Name,
		strings.Join(scopedDecl.BodyText, " "), scopedDecl.Comments, scopedDecl.PackageComments)

	expansions := make([]entity.Expansion, len(splits))
	for i, split := range splits {
		expansions[i] = entity.Expansion{From: split.Value, Values: amap.Expand(split.Value, scope, a.referenceText)}
	}

	return expansions
}

func (a amapExpander) ApplicableOn() string {
	return "samurai"
}
