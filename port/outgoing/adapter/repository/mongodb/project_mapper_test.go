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
		ID:        "715f17550be5f7222a815ff80966adaf",
		Status:    "done",
		Reference: "src-d/go-siva",
		CreatedAt: now,
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
	assert.Equal(t, "src-d/go-siva", dto.ProjecRef)
	assert.Equal(t, now, dto.CreatedAt)
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
	assert.Equal(t, int32(2), dto.SourceCode.FilesCount)
}

func TestToEntity_OnProjectMapper_ShouldReturnProjectEntity(t *testing.T) {
	now := time.Now()
	dto := projectDTO{
		ID:        "715f17550be5f7222a815ff80966adaf",
		Status:    "done",
		ProjecRef: "src-d/go-siva",
		CreatedAt: now,
		Metadata: metadataDTO{
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
		SourceCode: sourceCodeDTO{
			Hash:     "4ba248c1cf1003995d356f11935287b3e99decca",
			Location: "/tmp/repositories/github.com/src-d/go-siva",
			Files: []string{
				"common.go",
				"common_test.go",
			},
		},
	}

	pm := &projectMapper{}
	ent := pm.toEntity(dto)

	assert.Equal(t, "715f17550be5f7222a815ff80966adaf", ent.ID)
	assert.Equal(t, "done", ent.Status)
	assert.Equal(t, "src-d/go-siva", ent.Reference)
	assert.Equal(t, now, ent.CreatedAt)
	assert.Equal(t, "69565817", ent.Metadata.RemoteID)
	assert.Equal(t, "src-d", ent.Metadata.Owner)
	assert.Equal(t, "src-d/go-siva", ent.Metadata.Fullname)
	assert.Equal(t, "siva - seekable indexed verifiable archiver", ent.Metadata.Description)
	assert.Equal(t, "https://github.com/src-d/go-siva.git", ent.Metadata.CloneURL)
	assert.Equal(t, "master", ent.Metadata.DefaultBranch)
	assert.Equal(t, "mit", ent.Metadata.License)
	assert.Equal(t, now, *ent.Metadata.CreatedAt)
	assert.Equal(t, now, *ent.Metadata.UpdatedAt)
	assert.False(t, ent.Metadata.IsFork)
	assert.Equal(t, int32(102), ent.Metadata.Size)
	assert.Equal(t, int32(88), ent.Metadata.Stargazers)
	assert.Equal(t, int32(88), ent.Metadata.Watchers)
	assert.Equal(t, int32(16), ent.Metadata.Forks)
	assert.Equal(t, "4ba248c1cf1003995d356f11935287b3e99decca", ent.SourceCode.Hash)
	assert.Equal(t, "/tmp/repositories/github.com/src-d/go-siva", ent.SourceCode.Location)
	assert.ElementsMatch(t, []string{"common_test.go", "common.go"}, ent.SourceCode.Files)
}
