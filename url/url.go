package url

import (
	"regexp"
)

const (
	gitPattern   string = "^git@github.com:[A-Za-z0-9-]+/[\\w.-]+.git"
	httpsPattern string = "^https://github.com/[A-Za-z0-9-]+/[\\w-]+(.git)?$"
)

// IsValidGithubRepoURL checks if the given URL is valid. This validator validates against
// two types of GitHub URLs. It checks for `git@github.com` repository URLs and also for
// `https://github.com/` URLs. If the given string matches one of the types validator, the URL
// is considered valid.
func IsValidGithubRepoURL(url string) bool {
	gitMatcher := regexp.MustCompile(gitPattern)
	if gitMatcher.MatchString(url) {
		return true
	}

	httpsMatcher := regexp.MustCompile(httpsPattern)
	if httpsMatcher.MatchString(url) {
		return true
	}

	return false
}
