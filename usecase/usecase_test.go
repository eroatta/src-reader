package usecase_test

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/entity"
	"github.com/google/uuid"
)

// project repository mock
type projectRepositoryMock struct {
	project entity.Project
	getErr  error
	addErr  error
}

func (m projectRepositoryMock) Add(ctx context.Context, p entity.Project) error {
	return m.addErr
}

func (m projectRepositoryMock) Get(ctx context.Context, ID uuid.UUID) (entity.Project, error) {
	return m.project, m.getErr
}

func (m projectRepositoryMock) GetByReference(ctx context.Context, projectRef string) (entity.Project, error) {
	return m.project, m.getErr
}

// end project repository mock

// metadata repository mock
type metadataRepositoryMock struct {
	metadata entity.Metadata
	err      error
}

func (m metadataRepositoryMock) RetrieveMetadata(ctx context.Context, url string) (entity.Metadata, error) {
	return m.metadata, m.err
}

// end metadata repository mock

// source code repository mock
type sourceCodeRepositoryMock struct {
	sourceCode entity.SourceCode
	files      map[string][]byte
	err        error
}

func (m sourceCodeRepositoryMock) Clone(ctx context.Context, fullname string, url string) (entity.SourceCode, error) {
	return m.sourceCode, m.err
}

func (m sourceCodeRepositoryMock) Remove(ctx context.Context, location string) error {
	return m.err
}

func (m sourceCodeRepositoryMock) Read(ctx context.Context, location string, filename string) ([]byte, error) {
	if m.err != nil {
		return []byte{}, m.err
	}

	b, ok := m.files[filename]
	if !ok {
		return []byte{}, errors.New("not found")
	}

	return b, nil
}

// end source code repository mock

// identifier repository mock
type identifierRepositoryMock struct {
	idents []entity.Identifier
	err    error
}

func (i identifierRepositoryMock) Add(ctx context.Context, analysis entity.AnalysisResults, ident entity.Identifier) error {
	return i.err
}

func (i identifierRepositoryMock) FindAllByAnalysisID(ctx context.Context, analysisID uuid.UUID) ([]entity.Identifier, error) {
	return i.idents, i.err
}

func (i identifierRepositoryMock) FindAllByProjectAndFile(ctx context.Context, projectRef string, filename string) ([]entity.Identifier, error) {
	return i.idents, i.err
}

// end identifier repository mock

// analysis repository mock
type analysisRepositoryMock struct {
	analysisResults entity.AnalysisResults
	addErr          error
	getErr          error
}

func (a analysisRepositoryMock) Add(ctx context.Context, analysis entity.AnalysisResults) error {
	return a.addErr
}

func (a analysisRepositoryMock) GetByProjectID(ctx context.Context, projectID uuid.UUID) (entity.AnalysisResults, error) {
	return a.analysisResults, a.getErr
}

// end analysis repository mock

// insights repository mock
type insightsRepositoryMock struct {
	insights []entity.Insight
	addErr   error
	getErr   error
	delErr   error
}

func (i insightsRepositoryMock) AddAll(ctx context.Context, insights []entity.Insight) error {
	return i.addErr
}

func (i insightsRepositoryMock) GetByAnalysisID(ctx context.Context, analysisID uuid.UUID) ([]entity.Insight, error) {
	return i.insights, i.getErr
}

func (i insightsRepositoryMock) DeleteAllByAnalisysID(ctx context.Context, analysisID uuid.UUID) error {
	return i.delErr
}

// end insights repository mock
