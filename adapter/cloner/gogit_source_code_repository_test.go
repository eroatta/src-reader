package cloner

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/eroatta/src-reader/repository"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func TestNewGitCloneRepository_ShouldReturnNewInstance(t *testing.T) {
	sourceCodeRepository := NewGogitCloneRepository("/tmp/test", nil)

	assert.NotNil(t, sourceCodeRepository)
	assert.Equal(t, "/tmp/test", sourceCodeRepository.baseDir)
}

func TestClone_OnGogitCloneRepository_WhenUnableToCreateDestinationPath_ShouldReturnError(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "test-fail-create-folder-")
	if err != nil {
		assert.FailNow(t, "unexpected error creating temp folder", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile, err := ioutil.TempFile(tmpDir, "")
	if err != nil {
		assert.FailNow(t, "unexpected error creating temp folder", err)
	}
	defer os.Remove(tmpFile.Name())

	sourceCodeRepository := NewGogitCloneRepository(tmpDir, nil)
	existingFilename := strings.ReplaceAll(tmpFile.Name(), fmt.Sprintf("%s/", tmpDir), "")

	sourceCode, err := sourceCodeRepository.Clone(context.TODO(), existingFilename, "clone_url")

	assert.EqualError(t, err, repository.ErrSourceCodeUnableCreateDestination.Error())
	assert.Empty(t, sourceCode)
}

func TestClone_OnGogitCloneRepository_WhenUnableToCloneRemoteRepository_ShouldReturnError(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "test-fail-clone-")
	if err != nil {
		assert.FailNow(t, "unexpected error creating temp folder", err)
	}
	defer os.RemoveAll(tmpDir)

	errorFunc := func(ctx context.Context, path string, url string) (*git.Repository, error) {
		return nil, errors.New("oops! something failed...")
	}
	sourceCodeRepository := NewGogitCloneRepository(tmpDir, errorFunc)

	sourceCode, err := sourceCodeRepository.Clone(context.TODO(), "eroatta/testrepo", "clone_url")

	assert.EqualError(t, err, repository.ErrSourceCodeCloneRemoteRepository.Error())
	assert.Empty(t, sourceCode)
}

func TestClone_OnGogitCloneRepository_WhenUnableAccessHeadRef_ShouldReturnError(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "test-fail-accessing-head-")
	if err != nil {
		assert.FailNow(t, "unexpected error creating temp folder", err)
	}
	defer os.RemoveAll(tmpDir)

	// empty cloned repository, without HEAD ref
	rep, _ := git.Init(memory.NewStorage(), nil)
	clonerFunc := func(ctx context.Context, path string, url string) (*git.Repository, error) {
		return rep, nil
	}
	sourceCodeRepository := NewGogitCloneRepository(tmpDir, clonerFunc)

	sourceCode, err := sourceCodeRepository.Clone(context.TODO(), "eroatta/testrepo", "clone_url")

	assert.EqualError(t, err, repository.ErrSourceCodeUnableAccessMetadata.Error())
	assert.Empty(t, sourceCode)
}

func TestClone_OnGogitCloneRepository_WhenUnableRetrieveWorktree_ShouldReturnError(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "test-fail-reading-files-")
	if err != nil {
		assert.FailNow(t, "unexpected error creating temp folder", err)
	}
	defer os.RemoveAll(tmpDir)

	// empty cloned repository, including HEAD ref
	rep, _ := git.Init(memory.NewStorage(), nil)
	headRef := plumbing.NewHashReference(plumbing.HEAD, plumbing.NewHash("bc9968d75e48de59f0870ffb71f5e160bbbdcf52"))
	err = rep.Storer.SetReference(headRef)
	if err != nil {
		assert.FailNow(t, "unexpected error setting HEAD ref", err)
	}

	clonerFunc := func(ctx context.Context, path string, url string) (*git.Repository, error) {
		return rep, nil
	}
	sourceCodeRepository := NewGogitCloneRepository(tmpDir, clonerFunc)

	sourceCode, err := sourceCodeRepository.Clone(context.TODO(), "eroatta/testrepo", "clone_url")

	assert.EqualError(t, err, repository.ErrSourceCodeUnableAccessMetadata.Error())
	assert.Empty(t, sourceCode)
}

func TestClone_OnGogitCloneRepository_ShouldReturnSourceCode(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "test-clone-success")
	if err != nil {
		assert.FailNow(t, "unexpected error creating temp folder", err)
	}
	defer os.RemoveAll(tmpDir)

	// a cloned repository, including HEAD ref
	fs := memfs.New()
	files := []string{"main.go", "file.go", "file_test.go", "README.md", ".gitignore"}
	for _, name := range files {
		_, err = fs.Create(name)
		if err != nil {
			assert.FailNow(t, "unexpected error creating file", err)
		}
	}
	rep, _ := git.Init(memory.NewStorage(), fs)
	headRef := plumbing.NewHashReference(plumbing.HEAD, plumbing.NewHash("bc9968d75e48de59f0870ffb71f5e160bbbdcf52"))
	err = rep.Storer.SetReference(headRef)
	if err != nil {
		assert.FailNow(t, "unexpected error setting HEAD ref", err)
	}

	clonerFunc := func(ctx context.Context, path string, url string) (*git.Repository, error) {
		return rep, nil
	}
	sourceCodeRepository := NewGogitCloneRepository(tmpDir, clonerFunc)

	sourceCode, err := sourceCodeRepository.Clone(context.TODO(), "eroatta/testrepo", "clone_url")

	assert.NoError(t, err)
	assert.Equal(t, "bc9968d75e48de59f0870ffb71f5e160bbbdcf52", sourceCode.Hash)
	assert.Equal(t, fmt.Sprintf("%s/eroatta/testrepo", tmpDir), sourceCode.Location)
	assert.ElementsMatch(t, []string{"main.go", "file.go", "file_test.go", "README.md"}, sourceCode.Files)
}

func TestRemove_OnGogitCloneRepository_WithNonSharedBaseDir_ShouldReturnError(t *testing.T) {
	sourceCodeRepository := NewGogitCloneRepository("/tmp/mydir", nil)
	err := sourceCodeRepository.Remove(context.TODO(), "/tmp/another/dir")

	assert.EqualError(t, err, repository.ErrSourceCodeUnableToRemove.Error())
}

func TestRemove_OnGogitCloneRepository_WithExistingLocation_ShouldRemoveLocation(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		assert.FailNow(t, "unexpected error creating temp folder", err)
	}
	defer os.Remove(tmp)

	sourceCodeRepository := NewGogitCloneRepository(os.TempDir(), nil)
	err = sourceCodeRepository.Remove(context.TODO(), tmp)

	assert.NoError(t, err)
}
