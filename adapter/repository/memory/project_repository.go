package memory

import (
	"context"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
)

// InMemoryProjectRepository represents a In Memory database, focused on handling projects as memory elements.
type InMemoryProjectRepository struct {
	repos map[string]*entity.Project
}

// NewInMemoryProjectRepository created a repository.ProjectRepository backed up by memory storage.
func NewInMemoryProjectRepository() *InMemoryProjectRepository {
	return &InMemoryProjectRepository{
		repos: make(map[string]*entity.Project),
	}
}

// Add stores a Project entity into the underlying in memory storage.
func (r InMemoryProjectRepository) Add(ctx context.Context, project entity.Project) error {
	r.repos[project.URL] = &project
	return nil
}

// GetByURL finds an existing Project on the in memory storage, using the given URL as filter.
func (r InMemoryProjectRepository) GetByURL(ctx context.Context, url string) (entity.Project, error) {
	project, ok := r.repos[url]
	if !ok {
		return entity.Project{}, repository.ErrProjectNoResults
	}

	return *project, nil
}
