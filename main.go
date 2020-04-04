package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/eroatta/src-reader/adapter/cloner"
	"github.com/eroatta/src-reader/adapter/github"
	"github.com/eroatta/src-reader/adapter/persistence"
	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/extractor"
	"github.com/eroatta/src-reader/splitter"
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

func newGoodMain(url string) {
	// cloning step
	/*
		_, filesc, err := step.Clone(url, cl.New())
		if err != nil {
			log.Fatalf("Error reading repository %s: %v", url, err)
		}

		// parsing step
		parsedc := step.Parse(filesc)
		files := step.Merge(parsedc)

		// mining step
		countMiner := miner.NewCount()
		declarationMiner := miner.NewDeclaration(lists.Dicctionary)
		scopeMiner := miner.NewScope("fail")

		miningResults := step.Mine(files, countMiner, declarationMiner, scopeMiner)

		frequencyTable := samurai.NewFrequencyTable()

		countResults := miningResults[countMiner.Name()].(miner.Count)
		freq := countResults.Results().(map[string]int)
		for token, count := range freq {
			if len(token) == 1 {
				continue
			}
			frequencyTable.SetOccurrences(token, count)
		}

		declarationResults := miningResults[declarationMiner.Name()].(miner.Declaration)
		decls := declarationResults.Decls()
		for k := range decls {
			log.Println(fmt.Sprintf("Declaration: %s", k))
		}

		scopeResults := miningResults[scopeMiner.Name()].(miner.Scope)
		scopes := scopeResults.ScopedDeclarations()
		for k := range scopes {
			log.Println(fmt.Sprintf("Scope: %s", k))
		}

		// extraction step
		identc := step.Extract(files, extractor.New)

		// splitting step
		samuraiSplitter := splitter.NewSamurai(frequencyTable, frequencyTable)
		conservSplitter := splitter.NewConserv()
		greedySplitter := splitter.NewGreedy()

		splittedc := step.Split(identc, samuraiSplitter, conservSplitter, greedySplitter)

		// expansion step
		basicExpander := expander.NewBasic(declarationResults.Decls())
		// TODO: add reference text
		amapExpander := expander.NewAMAP(scopes, []string{})
		expandedc := step.Expand(splittedc, basicExpander, amapExpander)

		// storing step
		errors := step.Store(expandedc, storer.New())
		if len(errors) > 0 {
			log.Fatal("Something failed")
		}*/
}
