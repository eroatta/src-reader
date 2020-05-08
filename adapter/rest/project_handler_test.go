package rest_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/eroatta/src-reader/adapter/rest"
	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestPOST_OnProjectCreationHandler_WithoutBody_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer(nil, nil)

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
	router := rest.NewServer(nil, nil)

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
				"invalid field 'repository' with value null or empty"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnProjectCreationHandler_WithWrongDataType_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer(nil, nil)

	w := httptest.NewRecorder()
	body := `{
		"repository": 1
	}`
	req, _ := http.NewRequest("POST", "/projects", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"json: cannot unmarshal number into Go struct field postCreateProjectCommand.repository of type string"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnProjectCreationHandler_WithInvalidRepository_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer(nil, nil)

	w := httptest.NewRecorder()
	body := `{
		"repository": "./github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/projects", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"invalid field 'repository' with value ./github.com/eroatta/src-reader"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnProjectCreationHandler_WithInternalError_ShouldReturnHTTP500(t *testing.T) {
	router := rest.NewServer(mockCreateUsecase{
		p:   entity.Project{},
		err: errors.New("error accessing repository http://github.com/eroatta/src-reader"),
	}, nil)

	w := httptest.NewRecorder()
	body := `{
		"repository": "https://github.com/eroatta/src-reader"
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
	router := rest.NewServer(mockCreateUsecase{
		p: entity.Project{
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
		},
		err: nil,
	}, nil)

	w := httptest.NewRecorder()
	body := `{
		"repository": "https://github.com/src-d/go-siva"
	}`
	req, _ := http.NewRequest("POST", "/projects", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `
		{
			"id": "715f17550be5f7222a815ff80966adaf",
			"status": "done",
			"url": "https://github.com/src-d/go-siva",
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

type mockCreateUsecase struct {
	p   entity.Project
	err error
}

func (m mockCreateUsecase) Import(ctx context.Context, url string) (entity.Project, error) {
	return m.p, m.err
}
