package usecase_test

import (
	"context"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewGetProjectUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewGetProjectUsecase(nil)

	assert.NotNil(t, uc)
}

func TestProcess_OnGetProjectUsecase_WhenNoExistingProject_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		getErr: repository.ErrProjectNoResults,
	}

	uc := usecase.NewGetProjectUsecase(projectRepositoryMock)

	id, _ := uuid.NewUUID()
	project, err := uc.Process(context.TODO(), id)

	assert.Empty(t, project)
	assert.EqualError(t, err, usecase.ErrProjectNotFound.Error())
}

func TestProcess_OnGetProjectUsecase_WhenErrorRetrievingProject_ShouldReturnError(t *testing.T) {
	projectRepositoryMock := projectRepositoryMock{
		getErr: repository.ErrProjectUnexpected,
	}

	uc := usecase.NewGetProjectUsecase(projectRepositoryMock)

	id, _ := uuid.NewUUID()
	project, err := uc.Process(context.TODO(), id)

	assert.Empty(t, project)
	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnGetProjectUsecase_WhenExistingProject_ShouldReturnProject(t *testing.T) {
	id, _ := uuid.NewUUID()
	projectRepositoryMock := projectRepositoryMock{
		project: entity.Project{
			ID:        id.String(),
			Reference: "test/mytest",
			Status:    "done",
			Metadata: entity.Metadata{
				Fullname: "test/mytest",
			},
			SourceCode: entity.SourceCode{
				Hash: "196370e",
				Files: []string{
					"main.go",
				},
			},
		},
	}

	uc := usecase.NewGetProjectUsecase(projectRepositoryMock)

	project, err := uc.Process(context.TODO(), id)

	assert.NoError(t, err)
	assert.Equal(t, "done", project.Status)
	assert.Equal(t, "test/mytest", project.Reference)
	assert.Equal(t, "test/mytest", project.Metadata.Fullname)
	assert.Equal(t, "196370e", project.SourceCode.Hash)
	assert.Equal(t, 1, len(project.SourceCode.Files))
}
