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

func TestFilenames_RepositoryWith5Files_ShouldReturn5Names(t *testing.T) {
	// given a filesystem and a set of files on a repository
	fs := memfs.New()
	expected := []string{"main.go", "file.go", "file_test.go", "README.md", ".gitignore"}
	for _, name := range expected {
		fs.Create(name)
	}
	repository, err := git.Init(memory.NewStorage(), fs)

	// when we retrieve all the files from the repository
	got, err := Filenames(repository)

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 5, len(got), "number of files must be equal")
	assert.ElementsMatch(t, expected, got, "filenames don't match")
}

func TestFilenames_RepositoryWith2Files1Folder_ShouldReturn2Names(t *testing.T) {
	// given a filesystem and a set of files on a repository
	fs := memfs.New()
	expected := []string{"main.go", "file.go"}
	for _, name := range expected {
		fs.Create(name)
	}
	fs.MkdirAll("/ignored", 0666)
	repository, err := git.Init(memory.NewStorage(), fs)

	// when we retrieve all the files from the repository
	got, err := Filenames(repository)

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 2, len(got), "number of files must be equal")
	assert.ElementsMatch(t, expected, got, "filenames don't match")
}

func TestFilenames_RepositoryWith3FilesAnd2In1Folder_ShouldReturn3Names(t *testing.T) {
	// given a filesystem and a set of files on a repository
	fs := memfs.New()

	fs.Create("main.go")
	fs.MkdirAll("included", 0666)
	inFolder := []string{"file.go", "file_test.go"}
	for _, name := range inFolder {
		fs.Create("included/" + name)
	}
	fs.MkdirAll("ignored", 0666)
	repository, err := git.Init(memory.NewStorage(), fs)

	// when we retrieve all the files from the repository
	got, err := Filenames(repository)

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 3, len(got), "number of files must be equal")
	assert.ElementsMatch(t, []string{"main.go", "included/file.go", "included/file_test.go"}, got, "filenames don't match")
}

func TestFilenames_RepositoryWith1FileAnd2FilesOnSubSubFolder_ShouldReturn3Names(t *testing.T) {
	// given a filesystem and a set of files on a repository
	fs := memfs.New()

	fs.Create("main.go")
	fs.MkdirAll("sub/included", 0666)
	inFolder := []string{"file.go", "file_test.go"}
	for _, name := range inFolder {
		fs.Create(fs.Join("sub/included", name))
	}
	fs.MkdirAll("ignored", 0666)

	repository, err := git.Init(memory.NewStorage(), fs)

	// when we retrieve all the files from the repository
	got, err := Filenames(repository)

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 3, len(got), "number of files must be equal")
	assert.ElementsMatch(t, []string{"main.go", "sub/included/file.go", "sub/included/file_test.go"}, got, "filenames don't match")
}

func TestFilesnames_RepositoryNoFiles_ShouldReturn0Names(t *testing.T) {
	// given a filesystem but no files on a repository
	repository, err := git.Init(memory.NewStorage(), memfs.New())

	// when we retrieve all the files from the repository
	got, err := Filenames(repository)

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 0, len(got), "number of files must be equal")
	assert.ElementsMatch(t, []string{}, got, "filenames don't match")
}

func TestFilesInfo_RepositoryError_ShouldReturnAnError(t *testing.T) {
	// given a filesystem but no files on a repository
	repository, err := git.Init(memory.NewStorage(), nil)

	// when we retrieve all the files from the repository
	_, err = Filenames(repository)

	assert.NotNil(t, err, "error while retrieving the files from the repository should be nil")
}

func TestFile_RepositoryExistingFile_ShouldReturnAnArrayOfBytes(t *testing.T) {
	// given a filesystem and a set of files on a repository
	fs := memfs.New()
	file, err := fs.Create("main.go")
	_, err = file.Write([]byte("package main"))

	repository, err := git.Init(memory.NewStorage(), fs)

	// when we retrieve all the files from the repository
	got, err := File(repository, "main.go")

	assert.Nil(t, err, "error while reading an existing file should be nil")
	assert.Equal(t, got, []byte("package main"), "raw files should match")
}

func TestFile_RepositoryNoFile_ShouldReturnAnError(t *testing.T) {
	// given a filesystem and but no files on a repository
	repository, err := git.Init(memory.NewStorage(), memfs.New())

	// when we retrieve all the files from the repository
	got, err := File(repository, "any_file")

	assert.NotNil(t, err, "error while reading a non existing file shouldn't be nil")
	assert.EqualError(t, err, "file does not exist")
	assert.Nil(t, got, "raw file should be nil")
}
