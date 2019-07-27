package url

import (
	"regexp"
)

var (
	gitMatcher   = regexp.MustCompile("^git@github.com:[A-Za-z0-9-]+/[\\w.-]+.git")
	httpsMatcher = regexp.MustCompile("^https://github.com/[A-Za-z0-9-]+/[\\w-]+(.git)?$")
)

// IsValidGithubRepo checks if the given URL is valid. This validator validates against
// two types of GitHub URLs. It checks for `git@github.com` repository URLs and also for
// `https://github.com/` URLs. If the given string matches one of the types validator, the URL
// is considered valid.
func IsValidGithubRepo(url string) bool {
	if gitMatcher.MatchString(url) || httpsMatcher.MatchString(url) {
		return true
	}

	return false
}
