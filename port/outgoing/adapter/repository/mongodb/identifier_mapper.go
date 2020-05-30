package mongodb

import (
	"errors"
	"go/token"
	"strings"
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
	default:
		tokenString = "unknown"
	}

	return tokenString
}

func (im *identifierMapper) fromStringToToken(str string) token.Token {
	var tok token.Token
	switch str {
	case "func":
		tok = token.FUNC
	case "var":
		tok = token.VAR
	case "const":
		tok = token.CONST
	case "struct":
		tok = token.STRUCT
	case "interface":
		tok = token.INTERFACE
	default:
		tok = token.DEFAULT
	}

	return tok
}

// toDTO maps the entity for Identifier into a Data Transfer Object.
func (im *identifierMapper) toDTO(ent entity.Identifier, projectEnt entity.Project) identifierDTO {
	dto := identifierDTO{
		ID:              ent.ID,
		Package:         ent.Package,
		AbsolutePackage: ent.FullPackageName(),
		File:            ent.File,
		Position:        ent.Position,
		Name:            ent.Name,
		Type:            im.fromTokenToString(ent.Type),
		AnalysisID:      projectEnt.ID,
		ProjectRef:      projectEnt.Metadata.Fullname,
		CreatedAt:       time.Now(),
		Exported:        ent.Exported(),
		Normalization: normalizationDTO{
			Word:      ent.Normalization.Word,
			Algorithm: ent.Normalization.Algorithm,
			Score:     ent.Normalization.Score,
		},
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
				Order:              expansionEnt.Order,
				SplittingAlgorithm: expansionEnt.SplittingAlgorithm,
				From:               expansionEnt.From,
				Values:             expansionEnt.Values,
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

// toEntity maps the Data Transfer Object for Identifier into a domain entity.
func (im *identifierMapper) toEntity(dto identifierDTO) entity.Identifier {
	splits := make(map[string][]entity.Split, len(dto.Splits))
	for alg, split := range dto.Splits {
		items := make([]entity.Split, len(split))
		for i, splitDto := range split {
			items[i] = entity.Split{
				Order: splitDto.Order,
				Value: splitDto.Value,
			}
		}
		splits[alg] = items
	}

	expansions := make(map[string][]entity.Expansion, len(dto.Expansions))
	for alg, exp := range dto.Expansions {
		items := make([]entity.Expansion, len(exp))
		for i, expDto := range exp {
			items[i] = entity.Expansion{
				Order:              expDto.Order,
				SplittingAlgorithm: expDto.SplittingAlgorithm,
				From:               expDto.From,
				Values:             expDto.Values,
			}
		}
		expansions[alg] = items
	}

	var err error
	if dto.Error != "" {
		err = errors.New(dto.Error)
	}

	return entity.Identifier{
		ID:         dto.ID,
		Package:    dto.Package,
		File:       dto.File,
		Position:   dto.Position,
		Name:       dto.Name,
		Type:       im.fromStringToToken(dto.Type),
		Node:       nil,
		Splits:     splits,
		Expansions: expansions,
		Error:      err,
		Normalization: entity.Normalization{
			Word:      dto.Normalization.Word,
			Algorithm: dto.Normalization.Algorithm,
			Score:     dto.Normalization.Score,
		},
	}
}

// identifierDTO is the database representation for an Identifier.
type identifierDTO struct {
	ID               string                    `bson:"identifier_id"`
	Package          string                    `bson:"package"`
	AbsolutePackage  string                    `bson:"absolute_package"`
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
	Normalization    normalizationDTO          `bson:"normalization"`
}

// splitDTO is the database representation for an Identifier's Split results.
type splitDTO struct {
	Order int    `bson:"order"`
	Value string `bson:"value"`
}

// expansionDTO is the database representation for an Identifier's Expansion results.
type expansionDTO struct {
	Order              int      `bson:"order"`
	SplittingAlgorithm string   `bson:"splitting_algorithm"`
	From               string   `bson:"from"`
	Values             []string `bson:"values"`
}

// normalizationDTO is the database representation for an Identifer's Normalization results.
type normalizationDTO struct {
	Word      string  `bson:"word"`
	Algorithm string  `bson:"algorithm"`
	Score     float64 `bson:"score"`
}
