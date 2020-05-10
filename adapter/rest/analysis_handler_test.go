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
	"github.com/eroatta/src-reader/usecase/analyze"
	"github.com/stretchr/testify/assert"
)

func TestPOST_OnAnalysisCreationHandler_WithoutBody_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/analysis", nil)
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

func TestPOST_OnAnalysisCreationHandler_WithEmptyBody_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
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

func TestPOST_OnAnalysisCreationHandler_WithWrongDataType_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{
		"repository": 1
	}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"json: cannot unmarshal number into Go struct field createAnalysisCommand.repository of type string"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnAnalysisCreationHandler_WithInvalidRepository_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{
		"repository": "./github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
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

func TestPOST_OnAnalysisCreationHandler_WithNotFoundProject_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, mockAnalyzeUsecase{
		a:   entity.AnalysisResults{},
		err: analyze.ErrProjectNotFound,
	})

	w := httptest.NewRecorder()
	body := `{
		"repository": "https://github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"non-existing repository https://github.com/eroatta/src-reader"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnAnalysisCreationHandler_WithInternalError_ShouldReturnHTTP500(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, mockAnalyzeUsecase{
		a:   entity.AnalysisResults{},
		err: errors.New("error analyzing repository http://github.com/eroatta/src-reader"),
	})

	w := httptest.NewRecorder()
	body := `{
		"repository": "https://github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `
		{
			"name": "internal_error",
			"message": "internal server error",
			"details": [
				"error analyzing repository http://github.com/eroatta/src-reader"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnAnalysisCreationHandler_WithSuccess_ShouldReturnHTTP201(t *testing.T) {
	now := time.Date(2020, time.May, 5, 22, 0, 0, 0, time.UTC)
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, mockAnalyzeUsecase{
		a: entity.AnalysisResults{
			ID:                      "715f17550be5f7222a815ff80966adaf",
			ProjectName:             "src-d/go-siva",
			DateCreated:             now,
			PipelineMiners:          []string{"miner_1", "miner_2"},
			PipelineSplitters:       []string{"splitter_1", "splitter_2"},
			PipelineExpanders:       []string{"expander_1", "expander_2"},
			FilesTotal:              10,
			FilesValid:              8,
			FilesError:              2,
			FilesErrorSamples:       []string{"file_error"},
			IdentifiersTotal:        120,
			IdentifiersValid:        105,
			IdentifiersError:        15,
			IdentifiersErrorSamples: []string{"identifier_error"},
		},
		err: nil,
	})

	w := httptest.NewRecorder()
	body := `{
		"repository": "https://github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `
		{
			"id": "715f17550be5f7222a815ff80966adaf",
			"project_ref": "src-d/go-siva",
			"created_at": "2020-05-05T22:00:00Z",
			"miners": [
				"miner_1", 
				"miner_2"
			],
			"splitters": [
				"splitter_1",
				"splitter_2"
			],
			"expanders": [
				"expander_1",
				"expander_2"
			],
			"files_summary": {
				"total": 10,
				"valid": 8,
				"failed": 2,
				"error_samples": [
					"file_error"
				]
			},
			"identifiers_summary": {
				"total": 120,
				"valid": 105,
				"failed": 15,
				"error_samples": [
					"identifier_error"
				]
			}
		}`,
		w.Body.String())
}

type mockAnalyzeUsecase struct {
	a   entity.AnalysisResults
	err error
}

func (m mockAnalyzeUsecase) Analyze(ctx context.Context, url string) (entity.AnalysisResults, error) {
	return m.a, m.err
}
