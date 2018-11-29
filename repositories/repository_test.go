package repositories

import (
	"errors"
	"testing"

	"gopkg.in/src-d/go-git.v4/storage/memory"

	"gopkg.in/src-d/go-billy.v4/memfs"

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
	// given a filesystem and a set of files on a repository
	fs := memfs.New()
	expected := []string{"main.go", "file.go", "file_test.go", "README.md", ".gitignore"}
	for _, name := range expected {
		fs.Create(name)
	}
	repository, err := git.Init(memory.NewStorage(), fs)

	// when we retrieve all the files from the repository
	files, err := FilesInfo(repository)

	got := make([]string, 0)
	for _, file := range files {
		got = append(got, file.Name())
	}

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 5, len(files), "number of files must be equal")
	assert.ElementsMatch(t, expected, got, "filenames don't match")
}

func TestFilesInfo_RepositoryWith2Files1Folder_ShouldReturnAnArrayOfFileInfoWith2Elements(t *testing.T) {
	// given a filesystem and a set of files on a repository
	fs := memfs.New()
	expected := []string{"main.go", "file.go"}
	for _, name := range expected {
		fs.Create(name)
	}
	fs.MkdirAll("/ignored", 0666)
	repository, err := git.Init(memory.NewStorage(), fs)

	// when we retrieve all the files from the repository
	files, err := FilesInfo(repository)

	got := make([]string, 0)
	for _, file := range files {
		got = append(got, file.Name())
	}

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 2, len(files), "number of files must be equal")
	assert.ElementsMatch(t, expected, got, "filenames don't match")
}

func TestFilesInfo_RepositoryWith3FilesAnd2In1Folder_ShouldReturnAnArrayOfFileInfoWith3Elements(t *testing.T) {
	// given a filesystem and a set of files on a repository
	fs := memfs.New()

	fs.Create("main.go")
	fs.MkdirAll("/included", 0666)
	fs.Chroot("/included")
	inFolder := []string{"file.go", "file_test.go"}
	for _, name := range inFolder {
		fs.Create(name)
	}
	fs.MkdirAll("/ignored", 0666)
	repository, err := git.Init(memory.NewStorage(), fs)

	// when we retrieve all the files from the repository
	files, err := FilesInfo(repository)

	got := make([]string, 0)
	for _, file := range files {
		got = append(got, file.Name())
	}

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 3, len(files), "number of files must be equal")
	assert.ElementsMatch(t, []string{"main.go", "file.go", "file_test.go"}, got, "filenames don't match")
}

func TestFilesInfo_RepositoryNoFiles_ShouldReturnAnEmptyArray(t *testing.T) {
	// given a filesystem but no files on a repository
	repository, err := git.Init(memory.NewStorage(), memfs.New())

	// when we retrieve all the files from the repository
	files, err := FilesInfo(repository)

	got := make([]string, 0)
	for _, file := range files {
		got = append(got, file.Name())
	}

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 0, len(files), "number of files must be equal")
	assert.ElementsMatch(t, []string{}, got, "filenames don't match")
}

func TestFilesInfo_RepositoryError_ShouldReturnAnError(t *testing.T) {
	// given a filesystem but no files on a repository
	repository, err := git.Init(memory.NewStorage(), nil)

	// when we retrieve all the files from the repository
	files, err := FilesInfo(repository)

	got := make([]string, 0)
	for _, file := range files {
		got = append(got, file.Name())
	}

	assert.NotNil(t, err, "error while retrieving the files from the repository should be nil")
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
