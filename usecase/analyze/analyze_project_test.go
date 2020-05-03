package analyze_test

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"testing"

	"github.com/eroatta/src-reader/adapter/algorithm/splitter"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase/analyze"
	"github.com/stretchr/testify/assert"
)

func TestNewAnalyzeProjectUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := analyze.NewAnalyzeProjectUsecase(nil, nil)

	assert.Empty(t, uc)
}

func TestAnalyze_OnAnalyzeProjectUsecase_WhenFailingToReadFiles_ShouldReturnError(t *testing.T) {
	project := entity.Project{
		URL: "https://github.com/eroatta/test",
		SourceCode: entity.SourceCode{
			Hash:     "asdf1234asdf",
			Location: "/tmp/repositories/eroatta/test",
			Files:    []string{"main.go"},
		},
	}

	sourceCodeRepositoryMock := sourceCodeFileReaderMock{
		files: make(map[string][]byte),
		err:   repository.ErrSourceCodeUnableReadFile,
	}

	uc := analyze.NewAnalyzeProjectUsecase(sourceCodeRepositoryMock, nil)

	results, err := uc.Analyze(context.TODO(), project, &entity.AnalysisConfig{})

	assert.EqualError(t, err, analyze.ErrUnableToBuildASTs.Error())
	assert.Empty(t, results)
}

func TestAnalyze_OnAnalyzeProjectUsecase_WhenFailingToParseFiles_ShouldReturnError(t *testing.T) {
	project := entity.Project{
		URL: "https://github.com/eroatta/test",
		SourceCode: entity.SourceCode{
			Hash:     "asdf1234asdf",
			Location: "/tmp/repositories/eroatta/test",
			Files:    []string{"main.go"},
		},
	}

	sourceCodeRepositoryMock := sourceCodeFileReaderMock{
		files: map[string][]byte{
			"main.go": []byte("packa main"),
		},
		err: nil,
	}

	uc := analyze.NewAnalyzeProjectUsecase(sourceCodeRepositoryMock, nil)

	results, err := uc.Analyze(context.TODO(), project, &entity.AnalysisConfig{})

	assert.EqualError(t, err, analyze.ErrUnableToBuildASTs.Error())
	assert.Empty(t, results)
}

func TestAnalyze_OnAnalyzeProjectUsecase_WhenFailingToCreateSplitters_ShouldReturnError(t *testing.T) {
	project := entity.Project{
		URL: "https://github.com/eroatta/test",
		SourceCode: entity.SourceCode{
			Hash:     "asdf1234asdf",
			Location: "/tmp/repositories/eroatta/test",
			Files:    []string{"main.go"},
		},
	}

	sourceCodeRepositoryMock := sourceCodeFileReaderMock{
		files: map[string][]byte{
			"main.go": []byte("package main"),
		},
		err: nil,
	}

	uc := analyze.NewAnalyzeProjectUsecase(sourceCodeRepositoryMock, nil)

	results, err := uc.Analyze(context.TODO(), project, &entity.AnalysisConfig{
		Miners:    []string{},
		Splitters: []string{},
	})

	assert.EqualError(t, err, analyze.ErrUnableToCreateProcessors.Error())
	assert.Empty(t, results)
}

func TestAnalyze_OnAnalyzeProjectUsecase_WhenFailingToCreateExpanders_ShouldReturnError(t *testing.T) {
	project := entity.Project{
		URL: "https://github.com/eroatta/test",
		SourceCode: entity.SourceCode{
			Hash:     "asdf1234asdf",
			Location: "/tmp/repositories/eroatta/test",
			Files:    []string{"main.go"},
		},
	}

	sourceCodeRepositoryMock := sourceCodeFileReaderMock{
		files: map[string][]byte{
			"main.go": []byte("package main"),
		},
		err: nil,
	}

	uc := analyze.NewAnalyzeProjectUsecase(sourceCodeRepositoryMock, nil)

	results, err := uc.Analyze(context.TODO(), project, &entity.AnalysisConfig{
		Miners:                    []string{},
		Splitters:                 []string{"conserv"},
		SplittingAlgorithmFactory: splitter.NewSplitterFactory(),
		Expanders:                 []string{},
	})

	assert.EqualError(t, err, analyze.ErrUnableToCreateProcessors.Error())
	assert.Empty(t, results)
}

func TestAnalyze_OnAnalyzeProjectUsecase_WhenFailingToSaveIdentifiers_ShouldReturnError(t *testing.T) {
	project := entity.Project{
		URL: "https://github.com/eroatta/test",
		SourceCode: entity.SourceCode{
			Hash:     "asdf1234asdf",
			Location: "/tmp/repositories/eroatta/test",
			Files:    []string{"main.go"},
		},
	}

	sourceCodeRepositoryMock := sourceCodeFileReaderMock{
		files: map[string][]byte{
			"main.go": []byte("package main"),
		},
		err: nil,
	}

	identifierRepositoryMock := identifierRepositoryMock{
		err: repository.ErrIdentifierUnexpected,
	}

	uc := analyze.NewAnalyzeProjectUsecase(sourceCodeRepositoryMock, identifierRepositoryMock)

	results, err := uc.Analyze(context.TODO(), project, &entity.AnalysisConfig{
		Miners:                    []string{},
		ExtractorFactory:          newExtractorMock,
		Splitters:                 []string{"conserv"},
		SplittingAlgorithmFactory: splitter.NewSplitterFactory(),
		Expanders:                 []string{"mock"},
		ExpansionAlgorithmFactory: expanderAbstractFactoryMock{},
	})

	assert.EqualError(t, err, analyze.ErrUnableToSaveIdentifiers.Error())
	assert.Empty(t, results)
}

func TestAnalyze_OnAnalyzeProjectUsecase_WhenAnalyzingIdentifiers_ShouldReturnAnalysisResults(t *testing.T) {
	project := entity.Project{
		ID:  "asadfasa345asdfasdfa",
		URL: "https://github.com/eroatta/test",
		Metadata: entity.Metadata{
			Fullname: "eroatta/test",
		},
		SourceCode: entity.SourceCode{
			Hash:     "asdf1234asdf",
			Location: "/tmp/repositories/eroatta/test",
			Files:    []string{"main.go"},
		},
	}

	sourceCodeRepositoryMock := sourceCodeFileReaderMock{
		files: map[string][]byte{
			"main.go": []byte("package main"),
		},
		err: nil,
	}

	identifierRepositoryMock := identifierRepositoryMock{
		err: nil,
	}

	uc := analyze.NewAnalyzeProjectUsecase(sourceCodeRepositoryMock, identifierRepositoryMock)

	results, err := uc.Analyze(context.TODO(), project, &entity.AnalysisConfig{
		Miners:                    []string{},
		ExtractorFactory:          newExtractorMock,
		Splitters:                 []string{"conserv"},
		SplittingAlgorithmFactory: splitter.NewSplitterFactory(),
		Expanders:                 []string{"mock"},
		ExpansionAlgorithmFactory: expanderAbstractFactoryMock{},
	})

	assert.NoError(t, err)
	assert.Equal(t, "asadfasa345asdfasdfa", results.ID)
	assert.Equal(t, "eroatta/test", results.ProjectName)
	assert.Equal(t, 1, results.FilesTotal)
	assert.Equal(t, 1, results.FilesValid)
	assert.Equal(t, 0, results.FilesError)
	assert.Empty(t, results.FilesErrorSamples)
	assert.EqualValues(t, []string{}, results.PipelineMiners)
	assert.EqualValues(t, []string{"conserv"}, results.PipelineSplitters)
	assert.EqualValues(t, []string{"mock"}, results.PipelineExpanders)
	assert.Equal(t, 1, results.IdentifiersTotal)
	assert.Equal(t, 1, results.IdentifiersValid)
	assert.Equal(t, 0, results.IdentifiersError)
	assert.Empty(t, results.IdentifiersErrorSamples)
}

type sourceCodeFileReaderMock struct {
	files map[string][]byte
	err   error
}

func (m sourceCodeFileReaderMock) Clone(ctx context.Context, fullname string, cloneURL string) (entity.SourceCode, error) {
	return entity.SourceCode{}, errors.New("shouldn't be called")
}

func (m sourceCodeFileReaderMock) Remove(ctx context.Context, location string) error {
	return errors.New("shouldn't be called")
}

func (m sourceCodeFileReaderMock) Read(ctx context.Context, location string, filename string) ([]byte, error) {
	if m.err != nil {
		return []byte{}, m.err
	}

	b, ok := m.files[filename]
	if !ok {
		return []byte{}, errors.New("not found")
	}

	return b, nil
}

type expanderAbstractFactoryMock struct{}

func (e expanderAbstractFactoryMock) Get(name string) (entity.ExpanderFactory, error) {
	return expanderFactoryMock{}, nil
}

type expanderFactoryMock struct{}

func (e expanderFactoryMock) Make(miningResults map[string]entity.Miner) (entity.Expander, error) {
	return expanderMock{}, nil
}

type expanderMock struct{}

func (e expanderMock) Name() string {
	return "mock"
}

func (e expanderMock) ApplicableOn() string {
	return "conserv"
}

func (e expanderMock) Expand(ident entity.Identifier) []entity.Expansion {
	return []entity.Expansion{
		{From: ident.Name, Values: []string{fmt.Sprintf("%s-expanded", ident.Name)}}}
}

func newExtractorMock(filename string) entity.Extractor {
	return &extractorMock{
		idents: make([]entity.Identifier, 0),
	}
}

type extractorMock struct {
	idents []entity.Identifier
}

func (t *extractorMock) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch elem := node.(type) {
	case *ast.File:
		t.idents = append(t.idents, entity.Identifier{
			Name:       elem.Name.String(),
			Position:   elem.Pos(),
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		})
	}

	return t
}

func (t *extractorMock) Identifiers() []entity.Identifier {
	return t.idents
}

type identifierRepositoryMock struct {
	err error
}

func (i identifierRepositoryMock) Add(ctx context.Context, project entity.Project, ident entity.Identifier) error {
	return i.err
}
