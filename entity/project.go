package entity

import "time"

// TODO: change description
// Project defines a project under analysis.
type Project struct {
	Status   string
	Metadata Metadata
}

// Metadata holds the project information.
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
