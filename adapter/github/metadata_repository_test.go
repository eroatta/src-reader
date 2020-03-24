package github

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRESTMetadataRepository_ShouldReturnNewInstance(t *testing.T) {
	httpClient := &http.Client{}
	repository := NewRESTMetadataRepository(httpClient, "token")

	assert.NotNil(t, repository)
	assert.Equal(t, httpClient, repository.httpClient)
	assert.Equal(t, "token", repository.accessToken)
}

func TestRetrieveMetadata_OnRESTMetadataRepository_WhileInvalidToken_ShouldReturnError(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
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
