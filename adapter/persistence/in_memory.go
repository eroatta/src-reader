package persistence

import (
	"context"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
)

func NewInMemoryProjectRepository() *InMemoryProjectRepository {
	return &InMemoryProjectRepository{
		repos: make(map[string]*entity.Project),
	}
}

type InMemoryProjectRepository struct {
	repos map[string]*entity.Project
}

func (r InMemoryProjectRepository) Add(ctx context.Context, project entity.Project) error {
	r.repos[project.URL] = &project
	return nil
}

func (r InMemoryProjectRepository) GetByURL(ctx context.Context, url string) (entity.Project, error) {
	project, ok := r.repos[url]
	if !ok {
		return entity.Project{}, repository.ErrProjectNoResults
	}

	return *project, nil
}
