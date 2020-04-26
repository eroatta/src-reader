package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eroatta/src-reader/adapter/algorithm/expander"
	"github.com/eroatta/src-reader/adapter/algorithm/extractor"
	"github.com/eroatta/src-reader/adapter/algorithm/miner"
	"github.com/eroatta/src-reader/adapter/algorithm/splitter"
	"github.com/eroatta/src-reader/adapter/cloner"
	"github.com/eroatta/src-reader/adapter/github"
	"github.com/eroatta/src-reader/adapter/persistence"
	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/usecase/analyze"
	"github.com/eroatta/src-reader/usecase/create"
)

func main() {
	importProjectUsecase("https://github.com/src-d/go-siva")
}

func importProjectUsecase(url string) {
	projectRepository := persistence.NewInMemoryProjectRepository()
	remoteProjectRepository := github.NewRESTMetadataRepository(&http.Client{}, "https://api.github.com", os.Getenv("GITHUB_TOKEN"))
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
	output, err := os.OpenFile("csv_identifiers_repository.csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer output.Close()

	identifierRepository := persistence.NewCSVIdentifierRepository(output)
	analyzeUsecase := analyze.NewAnalyzeProjectUsecase(sourceCodeRepository, identifierRepository)

	analysisResults, err := analyzeUsecase.Analyze(context.TODO(), project, &entity.AnalysisConfig{
		Miners:                    []string{"wordcount", "scoped-declarations", "comments", "declarations", "global-frequency-table"},
		MinerAlgorithmFactory:     miner.NewMinerFactory(),
		ExtractorFactory:          extractor.New,
		Splitters:                 []string{"conserv", "greedy", "samurai"},
		SplittingAlgorithmFactory: splitter.NewSplitterFactory(),
		Expanders:                 []string{"noexp", "basic", "amap"},
		ExpansionAlgorithmFactory: expander.NewExpanderFactory(),
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Analysis::: Done")
	log.Println(fmt.Sprintf("Results -	Project -		ID: %s", analysisResults.ProjectID))
	log.Println(fmt.Sprintf("Results -	Project -		URL: %s", analysisResults.ProjectURL))
	log.Println(fmt.Sprintf("Results -	Files -			Total: %d", analysisResults.FilesTotal))
	log.Println(fmt.Sprintf("Results -	Files -			Valid: %d", analysisResults.FilesValid))
	log.Println(fmt.Sprintf("Results -	Files -			With Error: %d", analysisResults.FilesError))
	for _, sample := range analysisResults.FilesErrorSamples {
		log.Println(fmt.Sprintf("Results -	Files -			Error Sample: %s", sample))
	}
	log.Println(fmt.Sprintf("Results -	Pipeline -		Miners: %v", analysisResults.PipelineMiners))
	log.Println(fmt.Sprintf("Results -	Pipeline -		Splitters: %v", analysisResults.PipelineSplitters))
	log.Println(fmt.Sprintf("Results -	Pipeline -		Expanders: %v", analysisResults.PipelineExpanders))
	log.Println(fmt.Sprintf("Results -	Identifiers -		Total: %d", analysisResults.IdentifiersTotal))
	log.Println(fmt.Sprintf("Results -	Identifiers -		Valid: %d", analysisResults.IdentifiersValid))
	log.Println(fmt.Sprintf("Results -	Identifiers -		With Error: %d", analysisResults.IdentifiersError))
	for _, sample := range analysisResults.IdentifiersErrorSamples {
		log.Println(fmt.Sprintf("Results -	Identifiers -		Error Sample: %s", sample))
	}
}
