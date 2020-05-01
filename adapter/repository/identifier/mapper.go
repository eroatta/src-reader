package identifier

import (
	"go/token"

	"github.com/eroatta/src-reader/entity"
)

type identifierMapper struct{}

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

func (im *identifierMapper) toDTO(ent entity.Identifier, projectFullname string) identifierDTO {
	// setup direct mappings
	dto := identifierDTO{
		ID:         ent.ID,
		File:       ent.File,
		Position:   ent.Position,
		Name:       ent.Name,
		Type:       im.fromTokenToString(ent.Type),
		Parent:     ent.Parent,
		ParentPos:  ent.ParentPos,
		ProjectRef: projectFullname, // TODO: review approach
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
	ProjectRef string                    `bson:"project_ref_id"`
}

type splitDTO struct {
	Order int    `bson:"order"`
	Value string `bson:"value"`
}

type expansionDTO struct {
	From   string   `bson:"from"`
	Values []string `bson:"values"`
}
