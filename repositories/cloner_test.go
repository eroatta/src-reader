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

}

func TestClone_GithubClonerRepoNotFound_ShouldReturnAnError(t *testing.T) {

}

func TestGetFilesInfo_GithubClonerWith5Files_ShouldReturnAnArrayOfFileInfoWith5Elements(t *testing.T) {

}

func TestGetFilesInfo_GithubClonerWith2Files1Folder_ShouldReturnAnArrayOfFileInfoWith2Elements(t *testing.T) {

}

func TestGetFilesInfo_GithubClonerNoFiles_ShouldReturnAnEmptyArray(t *testing.T) {

}

func TestGetFilesInfo_GithubClonerError_ShouldReturnAnError(t *testing.T) {

}

func TestGetFile_GithubClonerEmptyFilename_ShouldReturnAnEmptyArray(t *testing.T) {

}

func TestGetFile_GithubClonerExistingFile_ShouldReturnAnArrayOfBytes(t *testing.T) {

}

func TestGetFile_GithubClonerNoFile_ShouldReturnAnError(t *testing.T) {

}
