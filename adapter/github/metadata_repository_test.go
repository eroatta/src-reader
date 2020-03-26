package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRESTMetadataRepository_ShouldReturnNewInstance(t *testing.T) {
	httpClient := &http.Client{}
	repository := NewRESTMetadataRepository(httpClient, "baseURL", "accessToken")

	assert.NotNil(t, repository)
	assert.Equal(t, httpClient, repository.httpClient)
	assert.Equal(t, "baseURL", repository.baseURL)
	assert.Equal(t, "accessToken", repository.accessToken)
}

func TestRetrieveMetadata_OnRESTMetadataRepository_WhileInvalidToken_ShouldReturnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		assert.Equal(t, "token invalid-token", accessToken)
		fmt.Println(accessToken)

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

	repository := NewRESTMetadataRepository(server.Client(), server.URL, "invalid-token")

	_, err := repository.RetrieveMetadata(context.TODO(), "owner/reponame")

	assert.Error(t, err)
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
