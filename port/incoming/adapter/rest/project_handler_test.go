package rest_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/port/incoming/adapter/rest"
	"github.com/eroatta/src-reader/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPOST_OnProjectCreationHandler_WithoutBody_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterCreateProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/projects", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"invalid request"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnProjectCreationHandler_WithEmptyBody_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterCreateProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{}`
	req, _ := http.NewRequest("POST", "/projects", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"invalid field 'reference' with value null or empty"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnProjectCreationHandler_WithWrongDataType_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterCreateProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{
		"reference": 1
	}`
	req, _ := http.NewRequest("POST", "/projects", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"json: cannot unmarshal number into Go struct field postCreateProjectCommand.reference of type string"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnProjectCreationHandler_WithInvalidRepository_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterCreateProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{
		"reference": "./github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/projects", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"invalid field 'reference' with value ./github.com/eroatta/src-reader"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnProjectCreationHandler_WithInternalError_ShouldReturnHTTP500(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterCreateProjectUsecase(router, mockCreateUsecase{
		project: entity.Project{},
		err:     errors.New("error accessing repository http://github.com/eroatta/src-reader"),
	})

	w := httptest.NewRecorder()
	body := `{
		"reference": "eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/projects", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `
		{
			"name": "internal_error",
			"message": "internal server error",
			"details": [
				"error accessing repository http://github.com/eroatta/src-reader"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnProjectCreationHandler_WithSuccess_ShouldReturnHTTP201(t *testing.T) {
	now := time.Date(2020, time.May, 5, 22, 0, 0, 0, time.UTC)
	router := rest.NewServer()
	rest.RegisterCreateProjectUsecase(router, mockCreateUsecase{
		project: entity.Project{
			ID:        uuid.MustParse("f9b76fde-c342-4328-8650-85da8f21e2be"),
			Status:    "done",
			Reference: "src-d/go-siva",
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
		},
		err: nil,
	})

	w := httptest.NewRecorder()
	body := `{
		"reference": "src-d/go-siva"
	}`
	req, _ := http.NewRequest("POST", "/projects", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `
		{
			"id": "f9b76fde-c342-4328-8650-85da8f21e2be",
			"status": "done",
			"reference": "src-d/go-siva",
			"metadata": {
				"remote_id": "69565817",
				"owner": "src-d",
				"fullname": "src-d/go-siva",
				"description": "siva - seekable indexed verifiable archiver",
				"clone_url": "https://github.com/src-d/go-siva.git",
				"branch": "master",
				"license": "mit",
				"created_at": "2020-05-05T22:00:00Z",
				"updated_at": "2020-05-05T22:00:00Z",
				"is_fork": false,
				"size": 102,
				"stargazers": 88,
				"watchers": 88,
				"forks": 16
			},
			"source_code": {
				"hash": "4ba248c1cf1003995d356f11935287b3e99decca",
				"location": "/tmp/repositories/github.com/src-d/go-siva",
				"files": [
					"common.go",
					"common_test.go"
				]
			}
		}`,
		w.Body.String())
}

func TestGET_OnProjectGetterHandler_WithInvalidID_ShouldReturnHTTP404(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGetProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects/asdfasdfs", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `
		{
			"name": "not_found",
			"message": "resource not found",
			"details": [
				"project with ID: asdfasdfs can't be found"
			]
		}`,
		w.Body.String())
}

func TestGET_OnProjectGetterHandler_WithNoExistingID_ShouldReturnHTTP404(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGetProjectUsecase(router, mockGetUsecase{
		err: usecase.ErrProjectNotFound,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects/6ba7b810-9dad-11d1-80b4-00c04fd430c8", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `
		{
			"name": "not_found",
			"message": "resource not found",
			"details": [
				"project with ID: 6ba7b810-9dad-11d1-80b4-00c04fd430c8 can't be found"
			]
		}`,
		w.Body.String())
}

func TestGET_OnProjectGetterHandler_WithInternalError_ShouldReturnHTTP500(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGetProjectUsecase(router, mockGetUsecase{
		err: usecase.ErrUnexpected,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects/6ba7b810-9dad-11d1-80b4-00c04fd430c8", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `
		{
			"name": "internal_error",
			"message": "internal server error",
			"details": [
				"error accessing project with ID: 6ba7b810-9dad-11d1-80b4-00c04fd430c8"
			]
		}`,
		w.Body.String())
}

func TestGET_OnProjectGetterHandler_WithSuccess_ShouldReturnHTTP200(t *testing.T) {
	now := time.Date(2020, time.May, 5, 22, 0, 0, 0, time.UTC)
	router := rest.NewServer()
	rest.RegisterGetProjectUsecase(router, mockGetUsecase{
		project: entity.Project{
			ID:        uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
			Status:    "done",
			Reference: "src-d/go-siva",
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
		},
		err: nil,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects/6ba7b810-9dad-11d1-80b4-00c04fd430c8", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `
		{
			"id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			"status": "done",
			"reference": "src-d/go-siva",
			"metadata": {
				"remote_id": "69565817",
				"owner": "src-d",
				"fullname": "src-d/go-siva",
				"description": "siva - seekable indexed verifiable archiver",
				"clone_url": "https://github.com/src-d/go-siva.git",
				"branch": "master",
				"license": "mit",
				"created_at": "2020-05-05T22:00:00Z",
				"updated_at": "2020-05-05T22:00:00Z",
				"is_fork": false,
				"size": 102,
				"stargazers": 88,
				"watchers": 88,
				"forks": 16
			},
			"source_code": {
				"hash": "4ba248c1cf1003995d356f11935287b3e99decca",
				"location": "/tmp/repositories/github.com/src-d/go-siva",
				"files": [
					"common.go",
					"common_test.go"
				]
			}
		}`,
		w.Body.String())
}

func TestDELETE_OnProjectDeleteHandler_WhenInvalidProjectID_ShouldReturn404(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterDeleteProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/projects/invalid-id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDELETE_OnProjectDeleteHandler_WhenErrorExecutingUsecase_ShouldReturn500(t *testing.T) {
	deleteProjectUsecaseMock := mockDeleteUsecase{
		err: usecase.ErrUnexpected,
	}

	router := rest.NewServer()
	rest.RegisterDeleteProjectUsecase(router, deleteProjectUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/projects/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `
		{
			"name": "internal_error",
			"message": "internal server error",
			"details": [
				"error deleting project with ID: ed2cd46a-4afd-4d49-a6ea-1c8d12d40134"
			]
		}`,
		w.Body.String())
}

func TestDELETE_OnProjectDeleteHandler_WhenDeletedProject_ShouldReturn204(t *testing.T) {
	deleteProjectUsecaseMock := mockDeleteUsecase{
		err: nil,
	}

	router := rest.NewServer()
	rest.RegisterDeleteProjectUsecase(router, deleteProjectUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/projects/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDELETE_OnProjectDeleteHandler_WhenAlreadyDeletedProject_ShouldReturn204(t *testing.T) {
	deleteDeleteUsecaseMock := mockDeleteUsecase{
		err: usecase.ErrProjectNotFound,
	}

	router := rest.NewServer()
	rest.RegisterDeleteProjectUsecase(router, deleteDeleteUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/projects/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

type mockCreateUsecase struct {
	project entity.Project
	err     error
}

func (m mockCreateUsecase) Process(ctx context.Context, projectRef string) (entity.Project, error) {
	return m.project, m.err
}

type mockGetUsecase struct {
	project entity.Project
	err     error
}

func (m mockGetUsecase) Process(ctx context.Context, ID uuid.UUID) (entity.Project, error) {
	return m.project, m.err
}

type mockDeleteUsecase struct {
	err error
}

func (m mockDeleteUsecase) Process(ctx context.Context, ID uuid.UUID) error {
	return m.err
}
