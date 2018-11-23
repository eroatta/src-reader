package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultGithubCloner_ShouldReturnNewCloner(t *testing.T) {
	cloner := DefaultGithubCloner()

	assert.Equal(t, "go-git", cloner.Name())
}
func TestClone_GithubCloner_ShouldCloneTheRepository(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestClone_GithubClonerRepoNotFound_ShouldReturnAnError(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestClone_GithubClonerConnectionError_ShouldReturnAnError(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}
