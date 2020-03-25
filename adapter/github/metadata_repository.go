package github

import (
	"context"
	"net/http"

	"github.com/eroatta/src-reader/entity"
)

func NewRESTMetadataRepository(httpClient *http.Client, baseURL string, accessToken string) *RESTMetadataRepository {
	return &RESTMetadataRepository{
		httpClient:  httpClient,
		baseURL:     baseURL,
		accessToken: accessToken,
	}
}

type RESTMetadataRepository struct {
	httpClient  *http.Client
	baseURL     string
	accessToken string
}

func (r RESTMetadataRepository) RetrieveMetadata(ctx context.Context, remoteRepository string) (entity.Metadata, error) {
	return entity.Metadata{}, nil
}
