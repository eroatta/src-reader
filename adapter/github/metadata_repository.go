package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"

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
	url := fmt.Sprintf("%s/repos/%s", r.baseURL, remoteRepository)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", fmt.Sprintf("token %s", r.accessToken))

	_, err := r.httpClient.Do(request)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("an error occurred while trying to GET the resource: %s", url))
		return entity.Metadata{}, repository.ErrUnexpected
	}

	return entity.Metadata{}, nil
}
