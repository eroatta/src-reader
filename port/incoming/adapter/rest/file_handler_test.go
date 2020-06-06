package rest_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eroatta/src-reader/port/incoming/adapter/rest"
	"github.com/eroatta/src-reader/usecase/file"
	"github.com/stretchr/testify/assert"
)

func TestGET_OnOriginalFileHandler_WithNoExistingProject_ShouldReturn404(t *testing.T) {
	originalFileUsecaseMock := originalFileUsecaseMock{
		t:                  t,
		expectedProjectRef: "eroatta/test",
		expectedFileRef:    "amap/amap.go",
		raw:                nil,
		err:                file.ErrProjectNotFound,
	}
	router := rest.NewServer()
	rest.RegisterOriginalFileUsecase(router, originalFileUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/files/originals/eroatta/test/amap/amap.go", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGET_OnOriginalFileHandler_WithNoExistingFile_ShouldReturn404(t *testing.T) {
	originalFileUsecaseMock := originalFileUsecaseMock{
		t:                  t,
		expectedProjectRef: "eroatta/test",
		expectedFileRef:    "amap/amap.go",
		raw:                nil,
		err:                file.ErrFileNotFound,
	}
	router := rest.NewServer()
	rest.RegisterOriginalFileUsecase(router, originalFileUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/files/originals/eroatta/test/amap/amap.go", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGET_OnOriginalFileHandler_WithErrorsWhileProcessing_ShouldReturn500(t *testing.T) {
	originalFileUsecaseMock := originalFileUsecaseMock{
		t:                  t,
		expectedProjectRef: "eroatta/test",
		expectedFileRef:    "amap/amap.go",
		raw:                nil,
		err:                file.ErrUnexpected,
	}
	router := rest.NewServer()
	rest.RegisterOriginalFileUsecase(router, originalFileUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/files/originals/eroatta/test/amap/amap.go", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Ooops! Something went wrong...", w.Body.String())
}

func TestGET_OnOriginalFileHandler_WithNoErrors_ShouldReturn200(t *testing.T) {
	content := []byte(`
		package main

		func main() {}
	`)
	originalFileUsecaseMock := originalFileUsecaseMock{
		t:                  t,
		expectedProjectRef: "eroatta/test",
		expectedFileRef:    "main.go",
		raw:                content,
		err:                nil,
	}
	router := rest.NewServer()
	rest.RegisterOriginalFileUsecase(router, originalFileUsecaseMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/files/originals/eroatta/test/main.go", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("content-type"))
	assert.Equal(t, string(content), w.Body.String())
}

type originalFileUsecaseMock struct {
	t                  *testing.T
	expectedProjectRef string
	expectedFileRef    string
	raw                []byte
	err                error
}

func (m originalFileUsecaseMock) Process(ctx context.Context, projectRef string, filename string) ([]byte, error) {
	assert.Equal(m.t, m.expectedProjectRef, projectRef)
	assert.Equal(m.t, m.expectedFileRef, filename)

	return m.raw, m.err
}
