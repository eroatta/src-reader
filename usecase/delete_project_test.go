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

func TestNewDeleteProjectUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewDeleteProjectUsecase(nil, nil, nil, nil)

	assert.Empty(t, uc)
}

func TestProcess_OnDeleteProjectUsecase_WhenErrorRetrievingAnalysis_ShouldReturnError(t *testing.T) {
	ar := analysisRepositoryMock{
		getErr: repository.ErrAnalysisUnexpected,
	}
	uc := usecase.NewDeleteProjectUsecase(nil, ar, nil, nil)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnDeleteProjectUsecase_WhenErrorDeletingAnalysis_ShouldReturnError(t *testing.T) {
	duc := deleteAnalysisUsecaseMock{
		err: usecase.ErrUnexpected,
	}
	ar := analysisRepositoryMock{
		getErr: repository.ErrAnalysisNoResults,
	}
	uc := usecase.NewDeleteProjectUsecase(duc, ar, nil, nil)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnDeleteProjectUsecase_WhenErrorRetrievingProject_ShouldReturnError(t *testing.T) {
	duc := deleteAnalysisUsecaseMock{
		err: nil,
	}
	ar := analysisRepositoryMock{
		getErr: repository.ErrAnalysisNoResults,
	}
	pr := projectRepositoryMock{
		getErr: repository.ErrProjectUnexpected,
	}
	uc := usecase.NewDeleteProjectUsecase(duc, ar, nil, pr)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnDeleteProjectUsecase_WhenErrorDeletingSourceCode_ShouldReturnError(t *testing.T) {
	duc := deleteAnalysisUsecaseMock{
		err: nil,
	}
	ar := analysisRepositoryMock{
		getErr: repository.ErrAnalysisNoResults,
	}
	scr := sourceCodeRepositoryMock{
		err: repository.ErrSourceCodeUnableToRemove,
	}
	pr := projectRepositoryMock{
		project: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "test/location",
			},
		},
		getErr: nil,
	}
	uc := usecase.NewDeleteProjectUsecase(duc, ar, scr, pr)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnDeleteProjectUsecase_WhenErrorDeletingProject_ShouldReturnError(t *testing.T) {
	duc := deleteAnalysisUsecaseMock{
		err: nil,
	}
	ar := analysisRepositoryMock{
		getErr: repository.ErrAnalysisNoResults,
	}
	scr := sourceCodeRepositoryMock{
		err: nil,
	}
	pr := projectRepositoryMock{
		project: entity.Project{
			SourceCode: entity.SourceCode{
				Location: "test/location",
			},
		},
		getErr: nil,
		delErr: repository.ErrProjectUnexpected,
	}
	uc := usecase.NewDeleteProjectUsecase(duc, ar, scr, pr)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrUnexpected.Error())
}

func TestProcess_OnDeleteProjectUsecase_WhenEverythingDeleted_ShouldReturnProjectNotFound(t *testing.T) {
	duc := deleteAnalysisUsecaseMock{
		err: nil,
	}
	ar := analysisRepositoryMock{
		getErr: repository.ErrAnalysisNoResults,
	}
	scr := sourceCodeRepositoryMock{
		err: nil,
	}
	pr := projectRepositoryMock{
		getErr: repository.ErrProjectNoResults,
		delErr: repository.ErrProjectNoResults,
	}
	uc := usecase.NewDeleteProjectUsecase(duc, ar, scr, pr)

	analysisID, _ := uuid.NewUUID()
	err := uc.Process(context.TODO(), analysisID)

	assert.EqualError(t, err, usecase.ErrProjectNotFound.Error())
}

type deleteAnalysisUsecaseMock struct {
	err error
}

func (m deleteAnalysisUsecaseMock) Process(ctx context.Context, analysisID uuid.UUID) error {
	return m.err
}
