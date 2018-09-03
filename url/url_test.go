package url_test

import (
	"testing"

	"github.com/eroatta/src-reader/url"
	"github.com/stretchr/testify/assert"
)

func TestGitRepositoryURLIsValid(t *testing.T) {
	cases := []string{
		"https://github.com/eroatta/dispersal",
		"https://github.com/eroatta/dispersal.git",
		"git@github.com:eroatta/dispersal.git",
	}

	for _, c := range cases {
		assert.True(t, url.IsValidGithubRepoURL(c), "Given URL should be valid")
	}
}

func TestGitRepoURLIsInvalid(t *testing.T) {
	cases := []string{
		"https://github.com/eroatta",
		"https://github.com/eroatta/dispersal/releases/tag/v0.0.1",
		"git@github.com:eroatta/dispersal",
		"github.com/eroatta/dispersal",
		"git@github.com/eroatta/.git",
	}

	for _, c := range cases {
		assert.False(t, url.IsValidGithubRepoURL(c), "Given URL should be invalid")
	}
}
