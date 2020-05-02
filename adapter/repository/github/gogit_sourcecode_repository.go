package github

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
)

// NewGogitSourceCodeRepository creates a new instance of SourceCodeRepository that clones source code
// using the src-d/go-git tool.
func NewGogitSourceCodeRepository(baseDir string, clonerFunc ClonerFunc) *GogitSourceCodeRepository {
	return &GogitSourceCodeRepository{
		baseDir:    baseDir,
		clonerFunc: clonerFunc,
	}
}

// ClonerFunc defines the interface for cloning a remote Git repository.
type ClonerFunc func(ctx context.Context, path string, url string) (*git.Repository, error)

// GogitSourceCodeRepository clones the remote source code repository into the local file system,
// using the src-d/go-git tool.
type GogitSourceCodeRepository struct {
	baseDir    string
	clonerFunc ClonerFunc
}

// PlainClonerFunc clones a remote GitHub repository using the src{d}/go-git client.
func PlainClonerFunc(ctx context.Context, path string, url string) (*git.Repository, error) {
	return git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL: url,
	})
}

func (r GogitSourceCodeRepository) Clone(ctx context.Context, fullname string, cloneURL string) (entity.SourceCode, error) {
	path := fmt.Sprintf("%s/%s", r.baseDir, fullname)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("failed to create directory %s", path))
		return entity.SourceCode{}, repository.ErrSourceCodeUnableCreateDestination
	}

	cloned, err := r.clonerFunc(ctx, path, cloneURL)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("failed to clone repository %s into %s", cloneURL, path))
		return entity.SourceCode{}, repository.ErrSourceCodeUnableCloneRemoteRepository
	}

	ref, err := cloned.Head()
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("failed to obtain ref to HEAD on cloned repository %s", cloneURL))
		return entity.SourceCode{}, repository.ErrSourceCodeUnableAccessMetadata
	}
	hash := ref.Hash().String()

	wt, err := cloned.Worktree()
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("failed to access worktree on cloned repository %s", cloneURL))
		return entity.SourceCode{}, repository.ErrSourceCodeUnableAccessMetadata
	}

	files, err := read(wt.Filesystem, "")
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("failed to get filenames on cloned repository %s", cloneURL))
		return entity.SourceCode{}, repository.ErrSourceCodeUnableAccessMetadata
	}

	return entity.SourceCode{
		Hash:     hash,
		Location: path,
		Files:    files,
	}, nil
}

func read(fs billy.Filesystem, rootDir string) ([]string, error) {
	files, err := fs.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		if file.IsDir() {
			subDirFilenames, err := read(fs, fs.Join(rootDir, file.Name()))
			if err != nil {
				return nil, err
			}

			names = append(names, subDirFilenames...)
		} else {
			names = append(names, fs.Join(rootDir, file.Name()))
		}
	}

	return names, nil
}

func (r GogitSourceCodeRepository) Remove(ctx context.Context, location string) error {
	if !strings.HasPrefix(location, r.baseDir) {
		return repository.ErrSourceCodeUnableToRemove
	}

	err := os.RemoveAll(location)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("unable to remove folder %s", location))
		return repository.ErrSourceCodeUnableToRemove
	}

	return nil
}

func (r GogitSourceCodeRepository) Read(ctx context.Context, location string, filename string) ([]byte, error) {
	if !strings.HasPrefix(location, r.baseDir) {
		return []byte{}, repository.ErrSourceCodeUnableReadFile
	}

	path := fmt.Sprintf("%s/%s", location, filename)
	file, err := os.Open(path)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("unable to access or open file %s", path))
		return []byte{}, repository.ErrSourceCodeUnableReadFile
	}

	raw, err := ioutil.ReadAll(file)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("unable read content on file %s", path))
		return []byte{}, repository.ErrSourceCodeUnableReadFile
	}

	return raw, nil
}
