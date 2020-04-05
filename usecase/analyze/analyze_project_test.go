package analyze_test

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"testing"

	"github.com/eroatta/src-reader/adapter/splitter"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/eroatta/src-reader/usecase/analyze"
	"github.com/stretchr/testify/assert"
)

func TestNewAnalyzeProjectUsecase_ShouldReturnNewInstance(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
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
		Miners:    []entity.Miner{},
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
		Miners:                    []entity.Miner{},
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
		err: repository.ErrUnexpected,
	}

	uc := analyze.NewAnalyzeProjectUsecase(sourceCodeRepositoryMock, identifierRepositoryMock)

	results, err := uc.Analyze(context.TODO(), project, &entity.AnalysisConfig{
		Miners:                    []entity.Miner{},
		ExtractorFactory:          newExtractorMock,
		Splitters:                 []string{"conserv"},
		SplittingAlgorithmFactory: splitter.NewSplitterFactory(),
		Expanders:                 []string{"mock"},
		ExpansionAlgorithmFactory: expanderAbstractFactoryMock{},
	})

	assert.EqualError(t, err, analyze.ErrUnableToSaveIdentifiers.Error())
	assert.Empty(t, results)
}

func TestAnalyze_OnAnalyzeProjectUsecase_WhenAnalyzingIdentifiers_ShouldResults(t *testing.T) {
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
		err: nil,
	}

	uc := analyze.NewAnalyzeProjectUsecase(sourceCodeRepositoryMock, identifierRepositoryMock)

	results, err := uc.Analyze(context.TODO(), project, &entity.AnalysisConfig{
		Miners:                    []entity.Miner{},
		ExtractorFactory:          newExtractorMock,
		Splitters:                 []string{"conserv"},
		SplittingAlgorithmFactory: splitter.NewSplitterFactory(),
		Expanders:                 []string{"mock"},
		ExpansionAlgorithmFactory: expanderAbstractFactoryMock{},
	})

	assert.NoError(t, err)
	assert.Empty(t, results)
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

func (e expanderFactoryMock) Make(staticInputs map[entity.InputType]interface{}, miningResults map[entity.MinerType]entity.Miner) (entity.Expander, error) {
	return expanderMock{}, nil
}

type expanderMock struct{}

func (e expanderMock) Name() string {
	return "mock"
}

func (e expanderMock) ApplicableOn() string {
	return "conserv"
}

func (e expanderMock) Expand(ident entity.Identifier) []string {
	return []string{fmt.Sprintf("%s-expanded", ident.Name)}
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
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
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
