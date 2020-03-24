package github

import "net/http"

func NewRESTMetadataRepository(httpClient *http.Client, accessToken string) *RESTMetadataRepository {
	return &RESTMetadataRepository{
		httpClient:  httpClient,
		accessToken: accessToken,
	}
}

type RESTMetadataRepository struct {
	httpClient  *http.Client
	accessToken string
}
