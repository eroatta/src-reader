package repositories

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"
)

func TestClone_ClonerOK_ShouldReturnRepository(t *testing.T) {
	repo, err := Clone(func(url string) (*git.Repository, error) {
		return &git.Repository{}, nil
	}, "git@github.com/test/case")

	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, repo, "cloned repository shouldn't be nil")
}
func TestClone_ClonerError_ShouldReturnAnError(t *testing.T) {
	repo, err := Clone(func(url string) (*git.Repository, error) {
		return nil, errors.New("connection error")
	}, "git@github.com/test/case")

	assert.Nil(t, repo, "cloned repository should be nil")
	assert.Equal(t, err.Error(), "Error cloning the remote repository")
}

// TODO add tests for GoGitClonerFunc

func TestFilesInfo_RepositoryWith5Files_ShouldReturnAnArrayOfFileInfoWith5Elements(t *testing.T) {

	expectedFilenames := []string{"main.go", "file.go", "file_test.go", "README.md", "license.txt"}

	repository := git.Repository{}

	files, err := FilesInfo(&repository)
	if err != nil {
		assert.Fail(t, "shouldn't get an error while retrieving the list of files")
	}

	filenames := make([]string, 0)
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	assert.Equal(t, 5, len(files), "number of files must be equal")
	assert.ElementsMatch(t, expectedFilenames, filenames, "filenames don't match")
}

func TestFilesInfo_RepositoryWith2Files1Folder_ShouldReturnAnArrayOfFileInfoWith2Elements(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestFilesInfo_RepositoryNoFiles_ShouldReturnAnEmptyArray(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestFilesInfo_RepositoryError_ShouldReturnAnError(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestFile_RepositoryEmptyFilename_ShouldReturnAnEmptyArrayOfBytes(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestFile_RepositoryExistingFile_ShouldReturnAnArrayOfBytes(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestFile_RepositoryNoFile_ShouldReturnAnError(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}
