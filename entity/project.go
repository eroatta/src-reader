package entity

import (
	"time"

	"github.com/google/uuid"
)

// Project represents a GitHub repository, which contains metadata about it
// and references to locally stored source code.
type Project struct {
	ID         uuid.UUID
	Status     string
	Reference  string
	CreatedAt  time.Time
	Metadata   Metadata
	SourceCode SourceCode
}

// Metadata holds the remote project information.
type Metadata struct {
	RemoteID      string
	Owner         string
	Fullname      string
	Description   string
	CloneURL      string
	DefaultBranch string
	License       string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	IsFork        bool
	Size          int32
	Stargazers    int32
	Watchers      int32
	Forks         int32
}

// SourceCode specifies the hash used for extracting the source code copy, its location
// and the associated list of included files.
type SourceCode struct {
	Hash     string
	Location string
	Files    []string
}
