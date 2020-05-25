package mongodb

import (
	"go/token"
	"strings"
	"time"
	"unicode"

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
		AnalysisID: projectEnt.ID,
		ProjectRef: projectEnt.Metadata.Fullname,
		CreatedAt:  time.Now(),
	}

	if unicode.IsUpper(rune(ent.Name[0])) && unicode.IsLetter(rune(ent.Name[0])) {
		dto.Exported = true
	}

	splits := make(map[string][]splitDTO, len(ent.Splits))
	joinedSplits := make(map[string]string, len(ent.Splits))
	for k, v := range ent.Splits {
		items := make([]splitDTO, len(v))
		words := make([]string, len(v))
		for i, splitEnt := range v {
			items[i] = splitDTO{
				Order: splitEnt.Order,
				Value: splitEnt.Value,
			}

			words[i] = splitEnt.Value
		}
		splits[k] = items
		joinedSplits[k] = strings.Join(words, "_")
	}
	dto.Splits = splits
	dto.JoinedSplits = joinedSplits

	expansions := make(map[string][]expansionDTO, len(ent.Expansions))
	joinedExpansions := make(map[string]string, len(ent.Expansions))
	for k, v := range ent.Expansions {
		items := make([]expansionDTO, len(v))
		words := make([]string, len(v))
		for i, expansionEnt := range v {
			items[i] = expansionDTO{
				From:   expansionEnt.From,
				Values: expansionEnt.Values,
			}

			words[i] = strings.Join(expansionEnt.Values, "|")
		}
		expansions[k] = items
		joinedExpansions[k] = strings.Join(words, "_")
	}
	dto.Expansions = expansions
	dto.JoinedExpansions = joinedExpansions

	return dto
}

// identifierDTO is the database representation for an Identifier.
type identifierDTO struct {
	ID               string                    `bson:"identifier_id"`
	File             string                    `bson:"file"`
	Position         token.Pos                 `bson:"position"`
	Name             string                    `bson:"name"`
	Type             string                    `bson:"type"`
	Splits           map[string][]splitDTO     `bson:"splits"`
	JoinedSplits     map[string]string         `bson:"joined_splits"`
	Expansions       map[string][]expansionDTO `bson:"expansions"`
	JoinedExpansions map[string]string         `bson:"joined_expansions"`
	Error            string                    `bson:"error_value,omitempty"`
	AnalysisID       string                    `bson:"analysis_id"`
	ProjectRef       string                    `bson:"project_ref"`
	CreatedAt        time.Time                 `bson:"created_at"`
	Exported         bool                      `bson:"is_exported"`
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
