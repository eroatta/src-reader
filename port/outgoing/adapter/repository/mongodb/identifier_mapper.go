package mongodb

import (
	"go/token"
	"time"

	"github.com/eroatta/src-reader/entity"
)

// identifierMapper maps an Identifier between its model and database representations.
type identifierMapper struct{}

// fromTokenToString transforms a token.Token value into a human-readable string.
func (im *identifierMapper) fromTokenToString(tok token.Token) string {
	var tokenString string
	switch tok {
	case token.FUNC:
		tokenString = "func"
	case token.VAR:
		tokenString = "var"
	case token.CONST:
		tokenString = "const"
	case token.STRUCT:
		tokenString = "struct"
	case token.INTERFACE:
		tokenString = "interface"
	case token.DEFINE:
		tokenString = "define"
	default:
		tokenString = "unknown"
	}

	return tokenString
}

// toDTO maps the entity for Identifier into a Data Transfer Object.
func (im *identifierMapper) toDTO(ent entity.Identifier, projectEnt entity.Project) identifierDTO {
	// setup direct mappings
	dto := identifierDTO{
		ID:         ent.ID,
		File:       ent.File,
		Position:   ent.Position,
		Name:       ent.Name,
		Type:       im.fromTokenToString(ent.Type),
		Parent:     ent.Parent,
		ParentPos:  ent.ParentPos,
		AnalysisID: projectEnt.ID,
		ProjectRef: projectEnt.Metadata.Fullname,
		CreatedAt:  time.Now(),
	}

	splits := make(map[string][]splitDTO, len(ent.Splits))
	for k, v := range ent.Splits {
		items := make([]splitDTO, len(v))
		for i, splitEnt := range v {
			items[i] = splitDTO{
				Order: splitEnt.Order,
				Value: splitEnt.Value,
			}
		}
		splits[k] = items
	}
	dto.Splits = splits

	expansions := make(map[string][]expansionDTO, len(ent.Expansions))
	for k, v := range ent.Expansions {
		items := make([]expansionDTO, len(v))
		for i, expansionEnt := range v {
			items[i] = expansionDTO{
				From:   expansionEnt.From,
				Values: expansionEnt.Values,
			}
		}
		expansions[k] = items
	}
	dto.Expansions = expansions

	return dto
}

// identifierDTO is the database representation for an Identifier.
type identifierDTO struct {
	ID         string                    `bson:"identifier_id"`
	File       string                    `bson:"file"`
	Position   token.Pos                 `bson:"position"`
	Name       string                    `bson:"name"`
	Type       string                    `bson:"type"`
	Parent     string                    `bson:"parent_id"`
	ParentPos  token.Pos                 `bson:"parent_position"`
	Splits     map[string][]splitDTO     `bson:"splits"`
	Expansions map[string][]expansionDTO `bson:"expansions"`
	Error      string                    `bson:"error_value,omitempty"`
	AnalysisID string                    `bson:"analysis_id"`
	ProjectRef string                    `bson:"project_ref"`
	CreatedAt  time.Time                 `bson:"created_at"`
}

// splitDTO is the database representation for an Identifier's Split results.
type splitDTO struct {
	Order int    `bson:"order"`
	Value string `bson:"value"`
}

// expansionDTO is the database representation for an Identifier's Expansion results.
type expansionDTO struct {
	From   string   `bson:"from"`
	Values []string `bson:"values"`
}
