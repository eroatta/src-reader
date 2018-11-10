package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClone_NoExistingGithubRepository_ShouldThrowError(t *testing.T) {
	_, err := Clone("http://github.com/mercadolibre/fury_credits-api")
	// "authentication required"
	// "http://githu.com/eroatta/repo" -> "invalid pkt-len found"
	assert.EqualError(t, err, "remote repository not found")
}

func TestClone_GithubRemoteError_ShouldThrowError(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}
func TestClone_GithubRepository_ShouldReturnRepository(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestRepository_FilesOnPlainRepo_ShouldReturnAllFiles(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestRepository_FilesOnRepoWithDirs_ShouldReturnAllFiles(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestReadDirAndSubDirs(t *testing.T) {
	//fs := memfs.New()

	//files, err := GetFiles()
	//assert.Equal(t, 4, len(files), "should match number of files")
}
