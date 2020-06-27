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
				"invalid field 'project_id' with value null or empty"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnAnalysisCreationHandler_WithWrongDataType_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{
		"project_id": 1
	}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"json: cannot unmarshal number into Go struct field createAnalysisCommand.project_id of type string"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnAnalysisCreationHandler_WithInvalidRepository_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{
		"project_id": "adsgadffadas"
	}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"invalid field 'project_id' with value adsgadffadas"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnAnalysisCreationHandler_WithNotFoundProject_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, mockAnalyzeUsecase{
		a:   entity.AnalysisResults{},
		err: usecase.ErrProjectNotFound,
	})

	w := httptest.NewRecorder()
	body := `{
		"project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"project with ID: 6ba7b810-9dad-11d1-80b4-00c04fd430c8 can't be found"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnAnalysisCreationHandler_WithAlreadyExistingAnalysis_ShouldReturnHTTP400(t *testing.T) {
	now := time.Date(2020, time.May, 5, 22, 0, 0, 0, time.UTC)
	router := rest.NewServer()
	rest.RegisterAnalyzeProjectUsecase(router, mockAnalyzeUsecase{
		a: entity.AnalysisResults{
			ID:                      uuid.MustParse("f17e675d-7823-4510-a04b-86e8c1f239ea"),
			ProjectName:             "src-d/go-siva",
			ProjectID:               uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
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
		err: usecase.ErrPreviousAnalysisFound,
	})

	w := httptest.NewRecorder()
	body := `{
		"project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"a previous analysis exists for project with ID: 6ba7b810-9dad-11d1-80b4-00c04fd430c8"
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
		"project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
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
			ID:                      uuid.MustParse("f17e675d-7823-4510-a04b-86e8c1f239ea"),
			ProjectName:             "src-d/go-siva",
			ProjectID:               uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
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
		"project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	}`
	req, _ := http.NewRequest("POST", "/analysis", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `
		{
			"id": "f17e675d-7823-4510-a04b-86e8c1f239ea",
			"created_at": "2020-05-05T22:00:00Z",
			"project_ref": "src-d/go-siva",
			"project_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
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

func TestDELETE_OnAnalysisDeleteHandler_WhenInvalidAnalysisID_ShouldReturn404(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterDeleteAnalysisUsecase(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/analysis/invalid-id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDELETE_OnAnalysisDeleteHandler_WhenErrorExecutingUsecase_ShouldReturn500(t *testing.T) {
	deleteAnalysisUsecaseMock := mockDeleteAnalysisUsecase{
		err: usecase.ErrUnexpected,
	}

	router := rest.NewServer()
	rest.RegisterDeleteAnalysisUsecase(router, deleteAnalysisUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/analysis/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `
		{
			"name": "internal_error",
			"message": "internal server error",
			"details": [
				"error deleting analysis with ID: ed2cd46a-4afd-4d49-a6ea-1c8d12d40134"
			]
		}`,
		w.Body.String())
}

func TestDELETE_OnAnalysisDeleteHandler_WhenDeletedAnalysis_ShouldReturn204(t *testing.T) {
	deleteAnalysisUsecaseMock := mockDeleteAnalysisUsecase{
		err: nil,
	}

	router := rest.NewServer()
	rest.RegisterDeleteAnalysisUsecase(router, deleteAnalysisUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/analysis/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDELETE_OnAnalysisDeleteHandler_WhenAlreadyDeletedAnalysis_ShouldReturn204(t *testing.T) {
	deleteAnalysisUsecaseMock := mockDeleteAnalysisUsecase{
		err: usecase.ErrAnalysisNotFound,
	}

	router := rest.NewServer()
	rest.RegisterDeleteAnalysisUsecase(router, deleteAnalysisUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/analysis/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

type mockAnalyzeUsecase struct {
	a   entity.AnalysisResults
	err error
}

func (m mockAnalyzeUsecase) Process(ctx context.Context, projectID uuid.UUID) (entity.AnalysisResults, error) {
	return m.a, m.err
}

type mockDeleteAnalysisUsecase struct {
	err error
}

func (m mockDeleteAnalysisUsecase) Process(ctx context.Context, analysisID uuid.UUID) error {
	return m.err
}
