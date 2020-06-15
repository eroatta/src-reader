package mongodb

import (
	"go/token"
	"testing"
	"time"

	"github.com/eroatta/src-reader/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFromTokenToString_OnIdentifierMapper_ShouldReturnTranslations(t *testing.T) {
	tests := []struct {
		name string
		tok  token.Token
		want string
	}{
		{"function_declaration", token.FUNC, "func"},
		{"variable_declaration", token.VAR, "var"},
		{"constant_declaration", token.CONST, "const"},
		{"struct_declaration", token.STRUCT, "struct"},
		{"interface_declaration", token.INTERFACE, "interface"},
		{"add_declaration", token.ADD, "unknown"},
	}

	im := &identifierMapper{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := im.fromTokenToString(tt.tok)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToDTO_OnIdentifierMapper_ShouldReturnIdentifierDTO(t *testing.T) {
	identifier := entity.Identifier{
		ID:       "filename:cmd/siva/impl/list.go+++pkg:impl+++declType:var+++name:defaultOutput",
		Package:  "impl",
		File:     "cmd/siva/impl/list.go",
		Position: token.Pos(194),
		Name:     "defaultOutput",
		Type:     token.VAR,
		Splits: map[string][]entity.Split{
			"conserv": {
				{Order: 1, Value: "default"},
				{Order: 2, Value: "output"},
			},
		},
		Expansions: map[string][]entity.Expansion{
			"noexp": {
				{Order: 1, SplittingAlgorithm: "conserv", From: "default", Values: []string{"default"}},
				{Order: 2, SplittingAlgorithm: "conserv", From: "output", Values: []string{"output"}},
			},
		},
		Normalization: entity.Normalization{
			Word:      "defaultOutput",
			Algorithm: "conserv+no_exp",
			Score:     0.99,
		},
	}
	analysis := entity.AnalysisResults{
		ID:          uuid.MustParse("f9b76fde-c342-4328-8650-85da8f21e2be"),
		ProjectName: "src-d/go-siva",
	}

	im := &identifierMapper{}
	dto := im.toDTO(identifier, analysis)

	assert.Equal(t, "filename:cmd/siva/impl/list.go+++pkg:impl+++declType:var+++name:defaultOutput", dto.ID)
	assert.Equal(t, "impl", dto.Package)
	assert.Equal(t, "cmd/siva/impl", dto.AbsolutePackage)
	assert.Equal(t, "cmd/siva/impl/list.go", dto.File)
	assert.Equal(t, token.Pos(194), dto.Position)
	assert.Equal(t, "defaultOutput", dto.Name)
	assert.Equal(t, "var", dto.Type)
	assert.Equal(t, 1, len(dto.Splits))
	assert.EqualValues(t, []splitDTO{
		{Order: 1, Value: "default"},
		{Order: 2, Value: "output"},
	}, dto.Splits["conserv"])
	assert.Equal(t, "default_output", dto.JoinedSplits["conserv"])
	assert.Equal(t, 1, len(dto.Expansions))
	assert.EqualValues(t, []expansionDTO{
		{Order: 1, SplittingAlgorithm: "conserv", From: "default", Values: []string{"default"}},
		{Order: 2, SplittingAlgorithm: "conserv", From: "output", Values: []string{"output"}},
	}, dto.Expansions["noexp"])
	assert.Equal(t, "default_output", dto.JoinedExpansions["noexp"])
	assert.Equal(t, "f9b76fde-c342-4328-8650-85da8f21e2be", dto.AnalysisID)
	assert.Equal(t, "src-d/go-siva", dto.ProjectRef)
	assert.Equal(t, time.Now().Format("2006-02-01"), dto.CreatedAt.Format("2006-02-01"))
	assert.False(t, dto.Exported)
	assert.Equal(t, "defaultOutput", dto.Normalization.Word)
	assert.Equal(t, "conserv+no_exp", dto.Normalization.Algorithm)
	assert.Equal(t, 0.99, dto.Normalization.Score)
}

func TestToEntity_OnIdentifierMapper_ShouldReturnIdentifierEntity(t *testing.T) {
	identifier := identifierDTO{
		ID:              "filename:cmd/siva/impl/list.go+++pkg:impl+++declType:var+++name:defaultOutput",
		Package:         "impl",
		AbsolutePackage: "cmd/siva/impl",
		File:            "cmd/siva/impl/list.go",
		Position:        194,
		Name:            "defaultOutput",
		Type:            "var",
		Splits: map[string][]splitDTO{
			"conserv": {
				{Order: 1, Value: "default"},
				{Order: 2, Value: "output"},
			},
		},
		JoinedSplits: map[string]string{
			"conserv": "default_output",
		},
		Expansions: map[string][]expansionDTO{
			"noexp": {
				{Order: 1, SplittingAlgorithm: "conserv", From: "default", Values: []string{"default"}},
				{Order: 2, SplittingAlgorithm: "conserv", From: "output", Values: []string{"output"}},
			},
		},
		JoinedExpansions: map[string]string{
			"noexp": "default_output",
		},
		AnalysisID: "715f17550be5f7222a815ff80966adaf",
		ProjectRef: "src-d/go-siva",
		CreatedAt:  time.Now(),
		Exported:   false,
		Normalization: normalizationDTO{
			Word:      "defaultOutput",
			Algorithm: "conserv+no_exp",
			Score:     0.99,
		},
	}

	im := &identifierMapper{}
	ent := im.toEntity(identifier)

	assert.Equal(t, "filename:cmd/siva/impl/list.go+++pkg:impl+++declType:var+++name:defaultOutput", ent.ID)
	assert.Equal(t, "impl", ent.Package)
	assert.Equal(t, "cmd/siva/impl", ent.FullPackageName())
	assert.Equal(t, "cmd/siva/impl/list.go", ent.File)
	assert.Equal(t, token.Pos(194), ent.Position)
	assert.Equal(t, "defaultOutput", ent.Name)
	assert.Equal(t, token.VAR, ent.Type)
	assert.Equal(t, 1, len(ent.Splits))
	assert.EqualValues(t, []entity.Split{
		{Order: 1, Value: "default"},
		{Order: 2, Value: "output"},
	}, ent.Splits["conserv"])
	assert.Equal(t, 1, len(ent.Expansions))
	assert.EqualValues(t, []entity.Expansion{
		{Order: 1, SplittingAlgorithm: "conserv", From: "default", Values: []string{"default"}},
		{Order: 2, SplittingAlgorithm: "conserv", From: "output", Values: []string{"output"}},
	}, ent.Expansions["noexp"])
	// assert.Equal(t, "715f17550be5f7222a815ff80966adaf", dto.AnalysisID)
	// assert.Equal(t, "src-d/go-siva", dto.ProjectRef)
	assert.False(t, ent.Exported())
	assert.Equal(t, "defaultOutput", ent.Normalization.Word)
	assert.Equal(t, "conserv+no_exp", ent.Normalization.Algorithm)
	assert.Equal(t, 0.99, ent.Normalization.Score)
}
