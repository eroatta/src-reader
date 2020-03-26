package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	response, err := r.httpClient.Do(request)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("an error occurred while trying to GET the resource: %s", url))
		return entity.Metadata{}, repository.ErrUnexpected
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.WithError(err).Error("an error occurred while reading response body")
		return entity.Metadata{}, repository.ErrUnexpected
	}

	if response.StatusCode != http.StatusOK {
		var errResponse errorResponse
		_ = json.Unmarshal(body, &errResponse)

		log.WithField("status_code", response.StatusCode).WithField("response_message", errResponse.Message).Error(fmt.Sprintf("an error occurred while trying to GET the resource: %s", url))
		return entity.Metadata{}, repository.ErrUnexpected
	}

	return entity.Metadata{}, nil
}

type errorResponse struct {
	Message string `json:"message"`
}
