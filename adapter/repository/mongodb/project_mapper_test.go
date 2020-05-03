package mongodb

import (
	"testing"
	"time"

	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestToDTO_OnProjectMapper_ShouldReturnProjectDTO(t *testing.T) {
	now := time.Now()
	entity := entity.Project{
		ID:     "715f17550be5f7222a815ff80966adaf",
		Status: "done",
		URL:    "https://github.com/src-d/go-siva",
		Metadata: entity.Metadata{
			RemoteID:      "69565817",
			Owner:         "src-d",
			Fullname:      "src-d/go-siva",
			Description:   "siva - seekable indexed verifiable archiver",
			CloneURL:      "https://github.com/src-d/go-siva.git",
			DefaultBranch: "master",
			License:       "mit",
			CreatedAt:     &now,
			UpdatedAt:     &now,
			IsFork:        false,
			Size:          102,
			Stargazers:    88,
			Watchers:      88,
			Forks:         16,
		},
		SourceCode: entity.SourceCode{
			Hash:     "4ba248c1cf1003995d356f11935287b3e99decca",
			Location: "/tmp/repositories/github.com/src-d/go-siva",
			Files: []string{
				"common.go",
				"common_test.go",
			},
		},
	}

	pm := &projectMapper{}
	dto := pm.toDTO(entity)

	assert.Equal(t, "715f17550be5f7222a815ff80966adaf", dto.ID)
	assert.Equal(t, "done", dto.Status)
	assert.Equal(t, "https://github.com/src-d/go-siva", dto.Url)
	assert.Equal(t, "69565817", dto.Metadata.RemoteID)
	assert.Equal(t, "src-d", dto.Metadata.Owner)
	assert.Equal(t, "src-d/go-siva", dto.Metadata.Fullname)
	assert.Equal(t, "siva - seekable indexed verifiable archiver", dto.Metadata.Description)
	assert.Equal(t, "https://github.com/src-d/go-siva.git", dto.Metadata.CloneURL)
	assert.Equal(t, "master", dto.Metadata.DefaultBranch)
	assert.Equal(t, "mit", dto.Metadata.License)
	assert.Equal(t, now, *dto.Metadata.CreatedAt)
	assert.Equal(t, now, *dto.Metadata.UpdatedAt)
	assert.False(t, dto.Metadata.IsFork)
	assert.Equal(t, int32(102), dto.Metadata.Size)
	assert.Equal(t, int32(88), dto.Metadata.Stargazers)
	assert.Equal(t, int32(88), dto.Metadata.Watchers)
	assert.Equal(t, int32(16), dto.Metadata.Forks)
	assert.Equal(t, "4ba248c1cf1003995d356f11935287b3e99decca", dto.SourceCode.Hash)
	assert.Equal(t, "/tmp/repositories/github.com/src-d/go-siva", dto.SourceCode.Location)
	assert.ElementsMatch(t, []string{"common_test.go", "common.go"}, dto.SourceCode.Files)
}
