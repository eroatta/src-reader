package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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
	// clean url in case of full remote address
	remoteRepository = strings.Replace(remoteRepository, "https://github.com/", "", -1)
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

	var okResponse successResponse
	err = json.Unmarshal(body, &okResponse)
	if err != nil {
		log.WithError(err).Error("an error occurred while parsing OK response body")
		return entity.Metadata{}, repository.ErrUnexpected
	}

	return entity.Metadata{
		RemoteID:      fmt.Sprintf("%d", okResponse.ID),
		Owner:         okResponse.Owner.Username,
		Fullname:      okResponse.Fullname,
		Description:   okResponse.Description,
		CloneURL:      okResponse.CloneURL,
		DefaultBranch: okResponse.Branch,
		License:       okResponse.License.ID,
		CreatedAt:     okResponse.CreatedAt,
		UpdatedAt:     okResponse.UpdatedAt,
		IsFork:        okResponse.IsFork,
		Size:          okResponse.Size,
		Stargazers:    okResponse.StargazersCount,
		Watchers:      okResponse.WatchersCount,
		Forks:         okResponse.ForksCount,
	}, nil
}

type errorResponse struct {
	Message string `json:"message"`
}

type successResponse struct {
	ID              int32      `json:"id"`
	Owner           owner      `json:"owner"`
	Fullname        string     `json:"full_name"`
	Description     string     `json:"description"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	CloneURL        string     `json:"clone_url"`
	Size            int32      `json:"size"`
	IsFork          bool       `json:"fork"`
	StargazersCount int32      `json:"stargazers_count"`
	WatchersCount   int32      `json:"watchers_count"`
	ForksCount      int32      `json:"forks"`
	License         license    `json:"license"`
	Branch          string     `json:"default_branch"`
}

type owner struct {
	Username string `json:"login"`
}

type license struct {
	ID string `json:"key"`
}
