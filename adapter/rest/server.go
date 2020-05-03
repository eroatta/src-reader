package rest

import (
	"github.com/eroatta/src-reader/usecase/analyze"
	"github.com/eroatta/src-reader/usecase/create"
	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", pingHandler)

	// r.POST("/projects", internal.createProject)
	// r.GET("/projects/$id", internal.getProject)
	// r.DELETE("/projects/$id", internal.deleteProject)
	// r.POST("/analysis", internal.createAnalysis)
	// r.GET("/analysis/$id", internal.getAnalysis)
	// r.DELETE("/analysis/$id", internal.deleteAnalysis)
	// r.GET("/analysis/$id/identifiers", internal.getIdentifiers)

	return r
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type server struct {
	createProjectUsecase  create.ImportProjectUsecase
	analyzeProjectUsecase analyze.AnalyzeProjectUsecase
}

func (s server) createAnalysis(ctx *gin.Context) {

}
