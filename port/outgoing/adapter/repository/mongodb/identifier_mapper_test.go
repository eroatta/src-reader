package mongodb

import (
	"go/token"
	"testing"
	"time"

	"github.com/eroatta/src-reader/entity"
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
		{"define_declaration", token.DEFINE, "define"},
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
		ID:        "filename:cmd/siva/impl/list.go+++pkg:impl+++declType:var+++name:defaultOutput",
		File:      "cmd/siva/impl/list.go",
		Position:  token.Pos(194),
		Name:      "defaultOutput",
		Type:      token.VAR,
		Parent:    "",
		ParentPos: token.NoPos,
		Splits: map[string][]entity.Split{
			"conserv": {
				{Order: 1, Value: "default"},
				{Order: 2, Value: "output"},
			},
		},
		Expansions: map[string][]entity.Expansion{
			"noexp": {
				{From: "default", Values: []string{"default"}},
				{From: "output", Values: []string{"output"}},
			},
		},
	}
	project := entity.Project{
		ID: "715f17550be5f7222a815ff80966adaf",
		Metadata: entity.Metadata{
			Fullname: "src-d/go-siva",
		},
	}

	im := &identifierMapper{}
	dto := im.toDTO(identifier, project)

	assert.Equal(t, "filename:cmd/siva/impl/list.go+++pkg:impl+++declType:var+++name:defaultOutput", dto.ID)
	assert.Equal(t, "cmd/siva/impl/list.go", dto.File)
	assert.Equal(t, token.Pos(194), dto.Position)
	assert.Equal(t, "defaultOutput", dto.Name)
	assert.Equal(t, "var", dto.Type)
	assert.Equal(t, "", dto.Parent)
	assert.Equal(t, token.NoPos, dto.ParentPos)
	assert.Equal(t, 1, len(dto.Splits))
	assert.EqualValues(t, []splitDTO{
		{Order: 1, Value: "default"},
		{Order: 2, Value: "output"},
	}, dto.Splits["conserv"])
	assert.Equal(t, 1, len(dto.Expansions))
	assert.EqualValues(t, []expansionDTO{
		{From: "default", Values: []string{"default"}},
		{From: "output", Values: []string{"output"}},
	}, dto.Expansions["noexp"])
	assert.Equal(t, "715f17550be5f7222a815ff80966adaf", dto.AnalysisID)
	assert.Equal(t, "src-d/go-siva", dto.ProjectRef)
	assert.Equal(t, time.Now().Format("2006-02-01"), dto.CreatedAt.Format("2006-02-01"))
	assert.False(t, dto.Exported)
}
