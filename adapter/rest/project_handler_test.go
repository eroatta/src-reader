package rest_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eroatta/src-reader/adapter/rest"
	"github.com/eroatta/src-reader/entity"
	"github.com/stretchr/testify/assert"
)

func TestPOST_OnProjectCreationHandler_WithoutBody_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer(nil)

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
	router := rest.NewServer(nil)

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
	router := rest.NewServer(nil)

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
	router := rest.NewServer(nil)

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
	router := rest.NewServer(mockUsecase{
		p:   entity.Project{},
		err: errors.New("error accessing repository http://github.com/eroatta/src-reader"),
	})

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
	router := rest.NewServer(mockUsecase{
		p: entity.Project{
			ID:     "13132dfadfasfasf",
			Status: "done",
			Metadata: entity.Metadata{
				Stargazers: 123,
			},
			SourceCode: entity.SourceCode{
				Hash:     "adfasdf9234adaf",
				Location: "/tmp/repositories/eroatta/src-reader",
				Files: []string{
					"main.go",
					"README.md",
				},
			},
		},
		err: nil,
	})

	w := httptest.NewRecorder()
	body := `{
		"repository": "https://github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/projects", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `
		{
			"id": "13132dfadfasfasf",
			"status": "done",
			"metadata": {
				"stargazers": 123
			},
			"source_code": {
				"hash": "adfasdf9234adaf",
				"location": "/tmp/repositories/eroatta/src-reader",
				"files": [
					"main.go",
					"README.md"
				]
			}
		}`,
		w.Body.String())
}

type mockUsecase struct {
	p   entity.Project
	err error
}

func (m mockUsecase) Import(ctx context.Context, url string) (entity.Project, error) {
	return m.p, m.err
}
