package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/port/incoming/adapter/rest"
	"github.com/eroatta/src-reader/port/outgoing/adapter/algorithm/expander"
	"github.com/eroatta/src-reader/port/outgoing/adapter/algorithm/extractor"
	"github.com/eroatta/src-reader/port/outgoing/adapter/algorithm/miner"
	"github.com/eroatta/src-reader/port/outgoing/adapter/algorithm/splitter"
	"github.com/eroatta/src-reader/port/outgoing/adapter/repository/github"
	"github.com/eroatta/src-reader/port/outgoing/adapter/repository/mongodb"
	"github.com/eroatta/src-reader/usecase/analyze"
	"github.com/eroatta/src-reader/usecase/create"
	log "github.com/sirupsen/logrus"
)

func main() {
	// create MongoDB client
	dbHost := os.Getenv("MONGODB_HOST")
	dbUsername := os.Getenv("MONGODB_USER")
	dbPassword := os.Getenv("MONGODB_PASSWORD")
	dbName := os.Getenv("MONGODB_DATABASE")
	clt, err := mongodb.NewMongoClient(fmt.Sprintf("mongodb://%s:%s@%s:27017/%s", dbUsername, dbPassword, dbHost, dbName))
	if err != nil {
		log.WithError(err).Fatal("Unable to start MongoDB client")
	}

	database := "reader"

	// create repositories based on MongoDB
	projectRepository := mongodb.NewMongoDBProjecRepository(clt, database)
	analysisRepository := mongodb.NewMongoDBAnalysisRepository(clt, database)
	identifierRepository := mongodb.NewMongoDBIdentifierRepository(clt, database)

	// create repositories based on Github
	githubToken := os.Getenv("GITHUB_TOKEN")
	sourceCodeFolder := "/tmp/repositories/github.com"
	remoteProjectRepository := github.NewRESTMetadataRepository(&http.Client{}, "https://api.github.com", githubToken)
	sourceCodeRepository := github.NewGogitSourceCodeRepository(sourceCodeFolder, github.PlainClonerFunc)

	// create supported use cases
	importProjectUsecase := create.NewImportProjectUsecase(projectRepository, remoteProjectRepository, sourceCodeRepository)
	analyzeProjectUsecase := analyze.NewAnalyzeProjectUsecase(projectRepository, sourceCodeRepository,
		identifierRepository, analysisRepository, defaultAnalysisConfig)

	// create REST API server and register use cases
	router := rest.NewServer()
	rest.RegisterCreateProjectUsecase(router, importProjectUsecase)
	rest.RegisterAnalyzeProjectUsecase(router, analyzeProjectUsecase)

	// start the server
	router.Run()
}

var defaultAnalysisConfig = &entity.AnalysisConfig{
	Miners:                    []string{"wordcount", "scoped-declarations", "comments", "declarations", "global-frequency-table"},
	MinerAlgorithmFactory:     miner.NewMinerFactory(),
	ExtractorFactory:          extractor.New,
	Splitters:                 []string{"conserv", "greedy", "samurai"},
	SplittingAlgorithmFactory: splitter.NewSplitterFactory(),
	Expanders:                 []string{"noexp", "basic", "amap"},
	ExpansionAlgorithmFactory: expander.NewExpanderFactory(),
}
