package mongodb

import (
	"time"

	"github.com/eroatta/src-reader/entity"
)

// projectMapper maps a Project between its model and database representations.
type projectMapper struct{}

// toDTO maps the entity for Project into a Data Transfer Object.
func (pm *projectMapper) toDTO(ent entity.Project) projectDTO {
	return projectDTO{
		ID:        ent.ID,
		Status:    ent.Status,
		Url:       ent.URL,
		CreatedAt: ent.CreatedAt,
		Metadata: metadataDTO{
			RemoteID:      ent.Metadata.RemoteID,
			Owner:         ent.Metadata.Owner,
			Fullname:      ent.Metadata.Fullname,
			Description:   ent.Metadata.Description,
			CloneURL:      ent.Metadata.CloneURL,
			DefaultBranch: ent.Metadata.DefaultBranch,
			License:       ent.Metadata.License,
			CreatedAt:     ent.Metadata.CreatedAt,
			UpdatedAt:     ent.Metadata.UpdatedAt,
			IsFork:        ent.Metadata.IsFork,
			Size:          ent.Metadata.Size,
			Stargazers:    ent.Metadata.Stargazers,
			Watchers:      ent.Metadata.Watchers,
			Forks:         ent.Metadata.Forks,
		},
		SourceCode: sourceCodeDTO{
			Hash:       ent.SourceCode.Hash,
			Location:   ent.SourceCode.Location,
			Files:      ent.SourceCode.Files,
			FilesCount: int32(len(ent.SourceCode.Files)),
		},
	}
}

// toEntity maps the Data Transfer Object for Project into a domain entity.
func (pm *projectMapper) toEntity(dto projectDTO) entity.Project {
	return entity.Project{
		ID:        dto.ID,
		Status:    dto.Status,
		URL:       dto.Url,
		CreatedAt: dto.CreatedAt,
		Metadata: entity.Metadata{
			RemoteID:      dto.Metadata.RemoteID,
			Owner:         dto.Metadata.Owner,
			Fullname:      dto.Metadata.Fullname,
			Description:   dto.Metadata.Description,
			CloneURL:      dto.Metadata.CloneURL,
			DefaultBranch: dto.Metadata.DefaultBranch,
			License:       dto.Metadata.License,
			CreatedAt:     dto.Metadata.CreatedAt,
			UpdatedAt:     dto.Metadata.UpdatedAt,
			IsFork:        dto.Metadata.IsFork,
			Size:          dto.Metadata.Size,
			Stargazers:    dto.Metadata.Stargazers,
			Watchers:      dto.Metadata.Watchers,
			Forks:         dto.Metadata.Forks,
		},
		SourceCode: entity.SourceCode{
			Hash:     dto.SourceCode.Hash,
			Location: dto.SourceCode.Location,
			Files:    dto.SourceCode.Files,
		},
	}
}

// projectDTO is the database representation for a Project.
type projectDTO struct {
	ID         string        `bson:"_id"`
	Status     string        `bson:"status"`
	Url        string        `bson:"url"`
	CreatedAt  time.Time     `bson:"created_at"`
	Metadata   metadataDTO   `bson:"metadata"`
	SourceCode sourceCodeDTO `bson:"source_code"`
}

// metadataDTO is the database representation for a Project's Metadata.
type metadataDTO struct {
	RemoteID      string     `bson:"remote_id"`
	Owner         string     `bson:"owner"`
	Fullname      string     `bson:"fullname"`
	Description   string     `bson:"description"`
	CloneURL      string     `bson:"clone_url"`
	DefaultBranch string     `bson:"branch"`
	License       string     `bson:"license"`
	CreatedAt     *time.Time `bson:"created_at"`
	UpdatedAt     *time.Time `bson:"updated_at"`
	IsFork        bool       `bson:"is_fork"`
	Size          int32      `bson:"size"`
	Stargazers    int32      `bson:"stargazers"`
	Watchers      int32      `bson:"watches"`
	Forks         int32      `bson:"forks"`
}

// sourceCodeDTO is the database representation for a Project's Source Code.
type sourceCodeDTO struct {
	Hash       string   `bson:"hash"`
	Location   string   `bson:"location"`
	Files      []string `bson:"files"`
	FilesCount int32    `bson:"files_count"`
}
