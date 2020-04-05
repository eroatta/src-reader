package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/eroatta/src-reader/adapter/cloner"
	"github.com/eroatta/src-reader/adapter/github"
	"github.com/eroatta/src-reader/adapter/persistence"
	"github.com/eroatta/src-reader/adapter/splitter"
	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/extractor"
	"github.com/eroatta/src-reader/usecase/analyze"
	"github.com/eroatta/src-reader/usecase/create"
)

func main() {
	importProjectUsecase("https://github.com/src-d/go-siva")
}

func importProjectUsecase(url string) {
	projectRepository := persistence.NewInMemoryProjectRepository()
	remoteProjectRepository := github.NewRESTMetadataRepository(&http.Client{}, "https://api.github.com", "")
	sourceCodeRepository := cloner.NewGogitCloneRepository("/tmp/repositories/github.com", cloner.PlainClonerFunc)

	uc := create.NewImportProjectUsecase(projectRepository, remoteProjectRepository, sourceCodeRepository)

	_, err := uc.Import(context.TODO(), url)
	if err != nil {
		log.Fatalln(err)
	}

	project, _ := projectRepository.GetByURL(context.TODO(), url)
	log.Println(project.Metadata.Owner)

	log.Println(project.SourceCode.Hash)
	log.Println(project.SourceCode.Location)
	for i, f := range project.SourceCode.Files {
		log.Println(fmt.Sprintf("File #%d: %s", i+1, f))
	}

	log.Println("Import::: Done")

	log.Println("Analysis::: Start")
	identiferRepository := identifierRepositoryMock{}
	analyzeUsecase := analyze.NewAnalyzeProjectUsecase(sourceCodeRepository, identiferRepository)

	_, err = analyzeUsecase.Analyze(context.TODO(), project, &entity.AnalysisConfig{
		Miners:                    make([]entity.Miner, 0),
		ExtractorFactory:          extractor.New,
		Splitters:                 []string{"conserv"},
		SplittingAlgorithmFactory: splitter.NewSplitterFactory(),
		Expanders:                 make([]string, 0),
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Analysis::: Done")
}

type identifierRepositoryMock struct {
}

func (i identifierRepositoryMock) Add(ctx context.Context, p entity.Project, ident code.Identifier) error {
	log.Println("Storing identifier...")
	for alg, splits := range ident.Splits {
		log.Println(fmt.Sprintf("%s \"%s\" Splitted into: %v by %s", ident.Type, ident.Name, splits, alg))
	}

	for alg, expans := range ident.Expansions {
		log.Println(fmt.Sprintf("%s \"%s\" Expanded into: %v by %s", ident.Type, ident.Name, expans, alg))
	}

	return nil
}
