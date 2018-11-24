package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDirAndSubDirs(t *testing.T) {
	//fs := memfs.New()

	//files, err := GetFiles()
	//assert.Equal(t, 4, len(files), "should match number of files")
}

func TestGetFilesInfo_GoGitRepositoryWith5Files_ShouldReturnAnArrayOfFileInfoWith5Elements(t *testing.T) {

	expectedFilenames := []string{"main.go", "file.go", "file_test.go", "README.md", "license.txt"}

	repository := GoGitRepository{}

	files, err := repository.GetFilesInfo()
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

func TestGetFilesInfo_GoGitRepositoryWith2Files1Folder_ShouldReturnAnArrayOfFileInfoWith2Elements(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestGetFilesInfo_GoGitRepositoryNoFiles_ShouldReturnAnEmptyArray(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestGetFilesInfo_GoGitRepositoryError_ShouldReturnAnError(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestGetFile_GoGitRepositoryEmptyFilename_ShouldReturnAnEmptyArray(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestGetFile_GoGitRepositoryExistingFile_ShouldReturnAnArrayOfBytes(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}

func TestGetFile_GoGitRepositoryNoFile_ShouldReturnAnError(t *testing.T) {
	assert.Fail(t, "unimplemented test")
}
