package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
		assert.Equal(t, "/repos/owner/reponame", r.RequestURI)

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

	assert.EqualError(t, err, repository.ErrMetadataUnexpected.Error())
	assert.Empty(t, metadata)
}

func TestRetrieveMetadata_OnRESTMetadataRepository_WhileNotFoundGitHubProject_ShouldReturnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		assert.Equal(t, "token valid-token", accessToken)
		assert.Equal(t, "/repos/owner/reponame", r.RequestURI)

		w.WriteHeader(http.StatusNotFound)
		body := `
			{
				"message": "Not Found",
				"documentation_url": "https://developer.github.com/v3/repos/#get"
		  	}
		`
		fmt.Fprintln(w, body)
	}))
	defer server.Close()

	metadataRepository := NewRESTMetadataRepository(server.Client(), server.URL, "valid-token")

	metadata, err := metadataRepository.RetrieveMetadata(context.TODO(), "owner/reponame")

	assert.EqualError(t, err, repository.ErrMetadataUnexpected.Error())
	assert.Empty(t, metadata)
}

func TestRetrieveMetadata_OnRESTMetadataRepository_WhileInternalError_ShouldReturnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		assert.Equal(t, "token valid-token", accessToken)
		assert.Equal(t, "/repos/owner/reponame", r.RequestURI)

		w.WriteHeader(http.StatusInternalServerError)
		body := `
			{
				"message": "Internal Server Error",
				"documentation_url": "https://developer.github.com/v3/repos/#get"
		  	}
		`
		fmt.Fprintln(w, body)
	}))
	defer server.Close()

	metadataRepository := NewRESTMetadataRepository(server.Client(), server.URL, "valid-token")

	metadata, err := metadataRepository.RetrieveMetadata(context.TODO(), "owner/reponame")

	assert.EqualError(t, err, repository.ErrMetadataUnexpected.Error())
	assert.Empty(t, metadata)
}

func TestRetrieveMetadata_OnRESTMetadataRepository_ShouldMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		assert.Equal(t, "token valid-token", accessToken)
		assert.Equal(t, "/repos/owner/reponame", r.RequestURI)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, bodyResponseOK)
	}))
	defer server.Close()

	metadataRepository := NewRESTMetadataRepository(server.Client(), server.URL, "valid-token")

	metadata, err := metadataRepository.RetrieveMetadata(context.TODO(), "https://github.com/owner/reponame")

	assert.NoError(t, err)
	assert.Equal(t, "223739110", metadata.RemoteID)
	assert.Equal(t, "eroatta", metadata.Owner)
	assert.Equal(t, "eroatta/freqtable", metadata.Fullname)
	assert.Equal(t, "Frequency Table builder, extracting and counting words from Go source code.", metadata.Description)
	assert.Equal(t, "https://github.com/eroatta/freqtable.git", metadata.CloneURL)
	assert.Equal(t, "master", metadata.DefaultBranch)
	assert.Equal(t, "mit", metadata.License)
	assert.Equal(t, "2019-11-24T12:16:58Z", metadata.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, "2020-02-16T20:07:09Z", metadata.UpdatedAt.Format(time.RFC3339))
	assert.False(t, metadata.IsFork)
	assert.Equal(t, int32(299), metadata.Size)
	assert.Equal(t, int32(0), metadata.Stargazers)
	assert.Equal(t, int32(0), metadata.Watchers)
	assert.Equal(t, int32(0), metadata.Forks)
}

var bodyResponseOK = `
{
	"id": 223739110,
	"node_id": "MDEwOlJlcG9zaXRvcnkyMjM3MzkxMTA=",
	"name": "freqtable",
	"full_name": "eroatta/freqtable",
	"private": false,
	"owner": {
	  "login": "eroatta",
	  "id": 5487897,
	  "node_id": "MDQ6VXNlcjU0ODc4OTc=",
	  "avatar_url": "https://avatars1.githubusercontent.com/u/5487897?v=4",
	  "gravatar_id": "",
	  "url": "https://api.github.com/users/eroatta",
	  "html_url": "https://github.com/eroatta",
	  "followers_url": "https://api.github.com/users/eroatta/followers",
	  "following_url": "https://api.github.com/users/eroatta/following{/other_user}",
	  "gists_url": "https://api.github.com/users/eroatta/gists{/gist_id}",
	  "starred_url": "https://api.github.com/users/eroatta/starred{/owner}{/repo}",
	  "subscriptions_url": "https://api.github.com/users/eroatta/subscriptions",
	  "organizations_url": "https://api.github.com/users/eroatta/orgs",
	  "repos_url": "https://api.github.com/users/eroatta/repos",
	  "events_url": "https://api.github.com/users/eroatta/events{/privacy}",
	  "received_events_url": "https://api.github.com/users/eroatta/received_events",
	  "type": "User",
	  "site_admin": false
	},
	"html_url": "https://github.com/eroatta/freqtable",
	"description": "Frequency Table builder, extracting and counting words from Go source code.",
	"fork": false,
	"url": "https://api.github.com/repos/eroatta/freqtable",
	"forks_url": "https://api.github.com/repos/eroatta/freqtable/forks",
	"keys_url": "https://api.github.com/repos/eroatta/freqtable/keys{/key_id}",
	"collaborators_url": "https://api.github.com/repos/eroatta/freqtable/collaborators{/collaborator}",
	"teams_url": "https://api.github.com/repos/eroatta/freqtable/teams",
	"hooks_url": "https://api.github.com/repos/eroatta/freqtable/hooks",
	"issue_events_url": "https://api.github.com/repos/eroatta/freqtable/issues/events{/number}",
	"events_url": "https://api.github.com/repos/eroatta/freqtable/events",
	"assignees_url": "https://api.github.com/repos/eroatta/freqtable/assignees{/user}",
	"branches_url": "https://api.github.com/repos/eroatta/freqtable/branches{/branch}",
	"tags_url": "https://api.github.com/repos/eroatta/freqtable/tags",
	"blobs_url": "https://api.github.com/repos/eroatta/freqtable/git/blobs{/sha}",
	"git_tags_url": "https://api.github.com/repos/eroatta/freqtable/git/tags{/sha}",
	"git_refs_url": "https://api.github.com/repos/eroatta/freqtable/git/refs{/sha}",
	"trees_url": "https://api.github.com/repos/eroatta/freqtable/git/trees{/sha}",
	"statuses_url": "https://api.github.com/repos/eroatta/freqtable/statuses/{sha}",
	"languages_url": "https://api.github.com/repos/eroatta/freqtable/languages",
	"stargazers_url": "https://api.github.com/repos/eroatta/freqtable/stargazers",
	"contributors_url": "https://api.github.com/repos/eroatta/freqtable/contributors",
	"subscribers_url": "https://api.github.com/repos/eroatta/freqtable/subscribers",
	"subscription_url": "https://api.github.com/repos/eroatta/freqtable/subscription",
	"commits_url": "https://api.github.com/repos/eroatta/freqtable/commits{/sha}",
	"git_commits_url": "https://api.github.com/repos/eroatta/freqtable/git/commits{/sha}",
	"comments_url": "https://api.github.com/repos/eroatta/freqtable/comments{/number}",
	"issue_comment_url": "https://api.github.com/repos/eroatta/freqtable/issues/comments{/number}",
	"contents_url": "https://api.github.com/repos/eroatta/freqtable/contents/{+path}",
	"compare_url": "https://api.github.com/repos/eroatta/freqtable/compare/{base}...{head}",
	"merges_url": "https://api.github.com/repos/eroatta/freqtable/merges",
	"archive_url": "https://api.github.com/repos/eroatta/freqtable/{archive_format}{/ref}",
	"downloads_url": "https://api.github.com/repos/eroatta/freqtable/downloads",
	"issues_url": "https://api.github.com/repos/eroatta/freqtable/issues{/number}",
	"pulls_url": "https://api.github.com/repos/eroatta/freqtable/pulls{/number}",
	"milestones_url": "https://api.github.com/repos/eroatta/freqtable/milestones{/number}",
	"notifications_url": "https://api.github.com/repos/eroatta/freqtable/notifications{?since,all,participating}",
	"labels_url": "https://api.github.com/repos/eroatta/freqtable/labels{/name}",
	"releases_url": "https://api.github.com/repos/eroatta/freqtable/releases{/id}",
	"deployments_url": "https://api.github.com/repos/eroatta/freqtable/deployments",
	"created_at": "2019-11-24T12:16:58Z",
	"updated_at": "2020-02-16T20:07:09Z",
	"pushed_at": "2020-02-16T20:07:07Z",
	"git_url": "git://github.com/eroatta/freqtable.git",
	"ssh_url": "git@github.com:eroatta/freqtable.git",
	"clone_url": "https://github.com/eroatta/freqtable.git",
	"svn_url": "https://github.com/eroatta/freqtable",
	"homepage": null,
	"size": 299,
	"stargazers_count": 0,
	"watchers_count": 0,
	"language": "Go",
	"has_issues": true,
	"has_projects": true,
	"has_downloads": true,
	"has_wiki": true,
	"has_pages": false,
	"forks_count": 0,
	"mirror_url": null,
	"archived": false,
	"disabled": false,
	"open_issues_count": 0,
	"license": {
	  "key": "mit",
	  "name": "MIT License",
	  "spdx_id": "MIT",
	  "url": "https://api.github.com/licenses/mit",
	  "node_id": "MDc6TGljZW5zZTEz"
	},
	"forks": 0,
	"open_issues": 0,
	"watchers": 0,
	"default_branch": "master",
	"temp_clone_token": null,
	"network_count": 0,
	"subscribers_count": 1
}
`
