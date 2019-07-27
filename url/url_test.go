package url_test

import (
	"testing"

	"github.com/eroatta/src-reader/url"
	"github.com/stretchr/testify/assert"
)

func TestIsValidGithubRepo_ShouldReturnTrueIfValidOrFalseOtherwise(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"regular_http_url", "https://github.com/eroatta/dispersal", true},
		{"regular_http_url_with_suffix", "https://github.com/eroatta/dispersal.git", true},
		{"regular_git_url", "git@github.com:eroatta/dispersal.git", true},
		{"missing_repository", "https://github.com/eroatta", false},
		{"unsupported_tags", "https://github.com/eroatta/dispersal/releases/tag/v0.0.1", false},
		{"git_url_without_suffix", "git@github.com:eroatta/dispersal", false},
		{"missing_protocol", "github.com/eroatta/dispersal", false},
		{"missing_repository_name", "git@github.com/eroatta/.git", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := url.IsValidGithubRepo(tt.url)

			assert.Equal(t, tt.want, got)
		})
	}
}
