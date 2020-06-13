package usecase_test

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
)

type projectRepositoryMock struct {
	project     entity.Project
	getByURLErr error
	addErr      error
}

func (m projectRepositoryMock) Add(ctx context.Context, p entity.Project) error {
	return m.addErr
}

func (m projectRepositoryMock) GetByURL(ctx context.Context, url string) (entity.Project, error) {
	return m.project, m.getByURLErr
}

type metadataRepositoryMock struct {
	metadata entity.Metadata
	err      error
}

func (m metadataRepositoryMock) RetrieveMetadata(ctx context.Context, url string) (entity.Metadata, error) {
	return m.metadata, m.err
}

type sourceCodeRepositoryMock struct {
	sourceCode entity.SourceCode
	err        error
}

func (m sourceCodeRepositoryMock) Clone(ctx context.Context, fullname string, url string) (entity.SourceCode, error) {
	return m.sourceCode, m.err
}

func (m sourceCodeRepositoryMock) Remove(ctx context.Context, location string) error {
	return m.err
}

func (m sourceCodeRepositoryMock) Read(ctx context.Context, location string, filename string) ([]byte, error) {
	return []byte{}, errors.New("shouldn't be called")
}

type identifierRepositoryMock struct {
	err error
}

func (i identifierRepositoryMock) Add(ctx context.Context, project entity.Project, ident entity.Identifier) error {
	return i.err
}

func (i identifierRepositoryMock) FindAllByProject(ctx context.Context, projectRef string) ([]entity.Identifier, error) {
	return []entity.Identifier{}, errors.New("shouldn't be called")
}

func (i identifierRepositoryMock) FindAllByProjectAndFile(ctx context.Context, projectRef string, filename string) ([]entity.Identifier, error) {
	return []entity.Identifier{}, errors.New("shouldn't be called")
}

type analysisRepositoryMock struct {
	err error
}

func (a analysisRepositoryMock) Add(ctx context.Context, analysis entity.AnalysisResults) error {
	return a.err
}
