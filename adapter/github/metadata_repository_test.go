package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eroatta/src-reader/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewRESTMetadataRepository_ShouldReturnNewInstance(t *testing.T) {
	httpClient := &http.Client{}
	metadataRepository := NewRESTMetadataRepository(httpClient, "baseURL", "accessToken")

	assert.NotNil(t, metadataRepository)
	assert.Equal(t, httpClient, metadataRepository.httpClient)
	assert.Equal(t, "baseURL", metadataRepository.baseURL)
	assert.Equal(t, "accessToken", metadataRepository.accessToken)
}

func TestRetrieveMetadata_OnRESTMetadataRepository_WhileInvalidToken_ShouldReturnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		assert.Equal(t, "token invalid-token", accessToken)

		w.WriteHeader(http.StatusUnauthorized)
		body := `
			{
				"message": "Bad credentials",
				"documentation_url": "https://developer.github.com/v3"
		  	}
		`
		fmt.Fprintln(w, body)
	}))
	defer server.Close()

	metadataRepository := NewRESTMetadataRepository(server.Client(), server.URL, "invalid-token")

	metadata, err := metadataRepository.RetrieveMetadata(context.TODO(), "owner/reponame")

	assert.EqualError(t, err, repository.ErrUnexpected.Error())
	assert.Empty(t, metadata)
}

func TestRetrieveMetadata_OnRESTMetadataRepository_WhileNotFoundGitHubProject_ShouldReturnError(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestRetrieveMetadata_OnRESTMetadataRepository_WhileInternalError_ShouldReturnError(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestRetrieveMetadata_OnRESTMetadataRepository_ShouldMetadata(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}
