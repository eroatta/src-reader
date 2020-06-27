package rest_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/port/incoming/adapter/rest"
	"github.com/eroatta/src-reader/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPOST_OnInsightsCreationHandler_WithoutBody_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/insights", nil)
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

func TestPOST_OnInsightsCreationHandler_WithEmptyBody_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"invalid field 'analysis_id' with value null or empty"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnInsightsCreationHandler_WithWrongDataType_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{
		"analysis_id": 1
	}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"json: cannot unmarshal number into Go struct field createInsightsCommand.analysis_id of type string"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnInsightsCreationHandler_WithInvalidRepository_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{
		"analysis_id": "./github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"invalid field 'analysis_id' with value ./github.com/eroatta/src-reader"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnInsightsCreationHandler_WithPreviousInsights_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, mockGainInsightsUsecase{
		ins: []entity.Insight{},
		err: usecase.ErrPreviousInsightsFound,
	})

	w := httptest.NewRecorder()
	body := `{
		"analysis_id": "a9f42bb8-92e6-4344-852b-2a9d8dd5b503"
	}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"previous insights exist for analysis with ID: a9f42bb8-92e6-4344-852b-2a9d8dd5b503"
			]
		}`,
		w.Body.String())
}
func TestPOST_OnInsightsCreationHandler_WithNotFoundIdentifiers_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, mockGainInsightsUsecase{
		ins: []entity.Insight{},
		err: usecase.ErrIdentifiersNotFound,
	})

	w := httptest.NewRecorder()
	body := `{
		"analysis_id": "a9f42bb8-92e6-4344-852b-2a9d8dd5b503"
	}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"non-existing identifiers for analysis ID a9f42bb8-92e6-4344-852b-2a9d8dd5b503"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnInsightsCreationHandler_WithInternalError_ShouldReturnHTTP500(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, mockGainInsightsUsecase{
		ins: []entity.Insight{},
		err: errors.New("error extracting insights from repository http://github.com/eroatta/src-reader"),
	})

	w := httptest.NewRecorder()
	body := `{
		"analysis_id": "a9f42bb8-92e6-4344-852b-2a9d8dd5b503"
	}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `
		{
			"name": "internal_error",
			"message": "internal server error",
			"details": [
				"error extracting insights from repository http://github.com/eroatta/src-reader"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnInsightsCreationHandler_WithSuccess_ShouldReturnHTTP201(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, mockGainInsightsUsecase{
		ins: []entity.Insight{
			{
				ProjectRef:       "eroatta/test",
				Package:          "main",
				TotalIdentifiers: 3,
				TotalExported:    1,
				TotalSplits: map[string]int{
					"conserv": 5,
				},
				TotalExpansions: map[string]int{
					"no_exp": 5,
				},
				TotalWeight: 2.267,
				Files: map[string]struct{}{
					"main.go":   {},
					"helper.go": {},
				},
			},
			{
				ProjectRef:       "eroatta/test",
				Package:          "main_test",
				TotalIdentifiers: 1,
				TotalExported:    1,
				TotalSplits: map[string]int{
					"conserv": 1,
				},
				TotalExpansions: map[string]int{
					"no_exp": 1,
				},
				TotalWeight: 1.0,
				Files: map[string]struct{}{
					"main_test.go": {},
				},
			},
		},
		err: nil,
	})

	w := httptest.NewRecorder()
	body := `{
		"analysis_id": "a9f42bb8-92e6-4344-852b-2a9d8dd5b503"
	}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `
		{
			"analysis_id": "a9f42bb8-92e6-4344-852b-2a9d8dd5b503",
			"project_ref": "eroatta/test",
			"identifiers": {
				"total": 4,
				"exported": 2
			},
			"accuracy": 0.81675,
			"packages": [
				{
					"name": "main",
					"accuracy": 0.7556666666666666,
					"identifiers": {
						"total": 3,
						"exported": 1
					},
					"files": [
						"helper.go",
						"main.go"
					]
				},
				{
					"name": "main_test",
					"accuracy": 1.00,
					"identifiers": {
						"total": 1,
						"exported": 1
					},
					"files": [
						"main_test.go"
					]
				}
			]
		}`,
		w.Body.String())
}

func TestGET_OnInsightsGetHandler_WithInvalidInsightsID_ShouldReturn404(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGetInsightsUsecase(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/insights/invalid-id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGET_OnInsightsGetHandler_WithNoExistingInsights_ShouldReturn404(t *testing.T) {
	retrieveInsightsMock := mockGetInsightsUsecase{
		ins: []entity.Insight{},
		err: usecase.ErrInsightsNotFound,
	}

	router := rest.NewServer()
	rest.RegisterGetInsightsUsecase(router, retrieveInsightsMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/insights/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGET_OnInsightsGetHandler_WithErrorReadingInsights_ShouldReturn500(t *testing.T) {
	retrieveInsightsMock := mockGetInsightsUsecase{
		ins: []entity.Insight{},
		err: usecase.ErrUnexpected,
	}

	router := rest.NewServer()
	rest.RegisterGetInsightsUsecase(router, retrieveInsightsMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/insights/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGET_OnInsightsGetHandler_WithExistingInsights_ShouldReturn200(t *testing.T) {
	retrieveInsightsMock := mockGetInsightsUsecase{
		ins: []entity.Insight{
			{
				ProjectRef:       "eroatta/test",
				Package:          "main",
				TotalIdentifiers: 3,
				TotalExported:    1,
				TotalSplits: map[string]int{
					"conserv": 5,
				},
				TotalExpansions: map[string]int{
					"no_exp": 5,
				},
				TotalWeight: 2.267,
				Files: map[string]struct{}{
					"main.go":   {},
					"helper.go": {},
				},
			},
			{
				ProjectRef:       "eroatta/test",
				Package:          "main_test",
				TotalIdentifiers: 1,
				TotalExported:    1,
				TotalSplits: map[string]int{
					"conserv": 1,
				},
				TotalExpansions: map[string]int{
					"no_exp": 1,
				},
				TotalWeight: 1.0,
				Files: map[string]struct{}{
					"main_test.go": {},
				},
			},
		},
		err: nil,
	}

	router := rest.NewServer()
	rest.RegisterGetInsightsUsecase(router, retrieveInsightsMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/insights/a9f42bb8-92e6-4344-852b-2a9d8dd5b503", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `
		{
			"analysis_id": "a9f42bb8-92e6-4344-852b-2a9d8dd5b503",
			"project_ref": "eroatta/test",
			"identifiers": {
				"total": 4,
				"exported": 2
			},
			"accuracy": 0.81675,
			"packages": [
				{
					"name": "main",
					"accuracy": 0.7556666666666666,
					"identifiers": {
						"total": 3,
						"exported": 1
					},
					"files": [
						"helper.go",
						"main.go"
					]
				},
				{
					"name": "main_test",
					"accuracy": 1.00,
					"identifiers": {
						"total": 1,
						"exported": 1
					},
					"files": [
						"main_test.go"
					]
				}
			]
		}`,
		w.Body.String())
}

func TestDELETE_OnInsightsDeleteHandler_WithInvalidAnalysisID_ShouldReturn404(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterDeleteInsightsUsecase(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/insights/invalid-id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDELETE_OnInsightsDeleteHandler_WhenNoExistingInsights_ShouldReturn204(t *testing.T) {
	deleteInsightsUsecaseMock := mockDeleteInsightsUsecase{
		err: usecase.ErrInsightsNotFound,
	}

	router := rest.NewServer()
	rest.RegisterDeleteInsightsUsecase(router, deleteInsightsUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/insights/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDELETE_OnInsightsDeleteHandler_WhenErrorDeletingInsights_ShouldReturn500(t *testing.T) {
	deleteInsightsUsecaseMock := mockDeleteInsightsUsecase{
		err: usecase.ErrUnexpected,
	}

	router := rest.NewServer()
	rest.RegisterDeleteInsightsUsecase(router, deleteInsightsUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/insights/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `
		{
			"name": "internal_error",
			"message": "internal server error",
			"details": [
				"error deleting insights for ID: ed2cd46a-4afd-4d49-a6ea-1c8d12d40134"
			]
		}`,
		w.Body.String())
}

func TestDELETE_OnInsightsDeleteHandler_WhenDeletedInsights_ShouldReturn204(t *testing.T) {
	deleteInsightsUsecaseMock := mockDeleteInsightsUsecase{
		err: nil,
	}

	router := rest.NewServer()
	rest.RegisterDeleteInsightsUsecase(router, deleteInsightsUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/insights/ed2cd46a-4afd-4d49-a6ea-1c8d12d40134", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

type mockGainInsightsUsecase struct {
	ins []entity.Insight
	err error
}

func (m mockGainInsightsUsecase) Process(ctx context.Context, analysisID uuid.UUID) ([]entity.Insight, error) {
	return m.ins, m.err
}

type mockGetInsightsUsecase struct {
	ins []entity.Insight
	err error
}

func (m mockGetInsightsUsecase) Process(ctx context.Context, insightsID uuid.UUID) ([]entity.Insight, error) {
	return m.ins, m.err
}

type mockDeleteInsightsUsecase struct {
	err error
}

func (m mockDeleteInsightsUsecase) Process(ctx context.Context, analysisID uuid.UUID) error {
	return m.err
}
