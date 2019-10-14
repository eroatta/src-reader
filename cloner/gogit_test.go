package cloner

import (
	"errors"
	"testing"

	"gopkg.in/src-d/go-git.v4/storage/memory"

	"gopkg.in/src-d/go-billy.v4/memfs"

	"github.com/eroatta/src-reader/code"
	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"
)

func TestClone_OnGoGitCloner_ShouldReturnRepository(t *testing.T) {
	clnr := goGitCloner{
		clonerFunc: func(url string) (*git.Repository, error) {
			return &git.Repository{}, nil
		},
	}
	repository, err := clnr.Clone("git@github.com/test/case")

	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, repository, "cloned repository shouldn't be nil")
}
func TestClone_OnGoGitClonerWithError_ShouldReturnAnError(t *testing.T) {
	clnr := goGitCloner{
		clonerFunc: func(url string) (*git.Repository, error) {
			return nil, errors.New("Connection error")
		},
	}
	repository, err := clnr.Clone("git@github.com/test/case")

	assert.Equal(t, code.Repository{}, repository, "cloned repository should be empty")
	assert.Equal(t, "Connection error", err.Error())
}

func TestFilenames_OnClonedRepositoryWith5Files_ShouldReturn5Names(t *testing.T) {
	// given a filesystem and a set of files on a repository
	fs := memfs.New()
	expected := []string{"main.go", "file.go", "file_test.go", "README.md", ".gitignore"}
	for _, name := range expected {
		fs.Create(name)
	}
	repository, err := git.Init(memory.NewStorage(), fs)
	clnr := goGitCloner{
		repository: repository,
	}

	// when we retrieve all the files from the repository
	got, err := clnr.Filenames()

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 5, len(got), "number of files must be equal")
	assert.ElementsMatch(t, expected, got, "filenames don't match")
}

func TestFilenames_OnClonedRepositoryWith2Files1Folder_ShouldReturn2Names(t *testing.T) {
	// given a filesystem and a set of files on a repository
	fs := memfs.New()
	expected := []string{"main.go", "file.go"}
	for _, name := range expected {
		fs.Create(name)
	}
	fs.MkdirAll("/ignored", 0666)
	repository, err := git.Init(memory.NewStorage(), fs)
	clnr := goGitCloner{
		repository: repository,
	}

	// when we retrieve all the files from the repository
	got, err := clnr.Filenames()

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 2, len(got), "number of files must be equal")
	assert.ElementsMatch(t, expected, got, "filenames don't match")
}

func TestFilenames_OnClonedRepositoryWith3FilesAnd2In1Folder_ShouldReturn3Names(t *testing.T) {
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
	clnr := goGitCloner{
		repository: repository,
	}

	// when we retrieve all the files from the repository
	got, err := clnr.Filenames()

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 3, len(got), "number of files must be equal")
	assert.ElementsMatch(t, []string{"main.go", "included/file.go", "included/file_test.go"}, got, "filenames don't match")
}

func TestFilenames_OnClonedRepositoryRepositoryWith1FileAnd2FilesOnSubSubFolder_ShouldReturn3Names(t *testing.T) {
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
	clnr := goGitCloner{
		repository: repository,
	}

	// when we retrieve all the files from the repository
	got, err := clnr.Filenames()

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 3, len(got), "number of files must be equal")
	assert.ElementsMatch(t, []string{"main.go", "sub/included/file.go", "sub/included/file_test.go"}, got, "filenames don't match")
}

func TestFilesnames_OnClonedRepositoryNoFiles_ShouldReturn0Names(t *testing.T) {
	// given a filesystem but no files on a repository
	repository, err := git.Init(memory.NewStorage(), memfs.New())
	clnr := goGitCloner{
		repository: repository,
	}

	// when we retrieve all the files from the repository
	got, err := clnr.Filenames()

	assert.Nil(t, err, "error while retrieving the files from the repository should be nil")
	assert.Equal(t, 0, len(got), "number of files must be equal")
	assert.ElementsMatch(t, []string{}, got, "filenames don't match")
}

func TestFilenames_OnClonedRepositoryError_ShouldReturnAnError(t *testing.T) {
	// given a filesystem but no files on a repository
	repository, err := git.Init(memory.NewStorage(), nil)
	clnr := goGitCloner{
		repository: repository,
	}

	// when we retrieve all the files from the repository
	_, err = clnr.Filenames()

	assert.NotNil(t, err, "error while retrieving the files from the repository should be nil")
}

func TestFile_OnClonedRepositoryExistingFile_ShouldReturnAnArrayOfBytes(t *testing.T) {
	// given a filesystem and a set of files on a repository
	fs := memfs.New()
	file, err := fs.Create("main.go")
	_, err = file.Write([]byte("package main"))

	repository, err := git.Init(memory.NewStorage(), fs)
	clnr := goGitCloner{
		repository: repository,
	}

	// when we retrieve all the files from the repository
	got, err := clnr.File("main.go")

	assert.Nil(t, err, "error while reading an existing file should be nil")
	assert.Equal(t, got, []byte("package main"), "raw files should match")
}

func TestFile_RepositoryNoFile_ShouldReturnAnError(t *testing.T) {
	// given a filesystem and but no files on a repository
	repository, err := git.Init(memory.NewStorage(), memfs.New())
	clnr := goGitCloner{
		repository: repository,
	}

	// when we retrieve all the files from the repository
	got, err := clnr.File("any_file")

	assert.NotNil(t, err, "error while reading a non existing file shouldn't be nil")
	assert.EqualError(t, err, "file does not exist")
	assert.Nil(t, got, "raw file should be nil")
}
