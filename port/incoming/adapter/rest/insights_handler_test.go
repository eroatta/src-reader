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
	"github.com/eroatta/src-reader/usecase/gain"
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
				"invalid field 'repository' with value null or empty"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnInsightsCreationHandler_WithWrongDataType_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{
		"repository": 1
	}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"json: cannot unmarshal number into Go struct field createInsightsCommand.repository of type string"
			]
		}`,
		w.Body.String())
}

func TestPOST_OnInsightsCreationHandler_WithInvalidRepository_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, nil)

	w := httptest.NewRecorder()
	body := `{
		"repository": "./github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
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

func TestPOST_OnInsightsCreationHandler_WithNotFoundIdentifiers_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer()
	rest.RegisterGainInsightsUsecase(router, mockGainInsightsUsecase{
		ins: []entity.Insight{},
		err: gain.ErrIdentifiersNotFound,
	})

	w := httptest.NewRecorder()
	body := `{
		"repository": "https://github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"name": "validation_error",
			"message": "missing or invalid data",
			"details": [
				"non-existing identifiers for https://github.com/eroatta/src-reader"
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
		"repository": "https://github.com/eroatta/src-reader"
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
		"repository": "https://github.com/eroatta/src-reader"
	}`
	req, _ := http.NewRequest("POST", "/insights", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `
		{
			"project_ref": "https://github.com/eroatta/src-reader",
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

type mockGainInsightsUsecase struct {
	ins []entity.Insight
	err error
}

func (m mockGainInsightsUsecase) Process(ctx context.Context, url string) ([]entity.Insight, error) {
	return m.ins, m.err
}
