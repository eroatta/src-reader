package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eroatta/src-reader/usecase/analyze"
	"github.com/eroatta/src-reader/usecase/create"
	"github.com/eroatta/token/conserv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// requestValidator represents a validator capable of analyzing the values of the incoming
// request bodies.
var requestValidator = validator.New()

func NewServer(ipUsecase create.ImportProjectUsecase) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", pingHandler)

	internal := server{
		createProjectUsecase: ipUsecase,
	}

	r.POST("/projects", internal.createProject)
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

type errorResponse struct {
	Name    string   `json:"name"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func newBadRequestResponse() errorResponse {
	return errorResponse{
		Name:    "validation_error",
		Message: "missing or invalid data",
		Details: make([]string, 0),
	}
}

func setBadRequestOnBindingResponse(ctx *gin.Context, err error) {
	errResponse := newBadRequestResponse()
	errResponse.Details = append(errResponse.Details, err.Error())

	ctx.JSON(http.StatusBadRequest, errResponse)
}

func setBadRequestOnValidationResponse(ctx *gin.Context, err error) {
	errResponse := newBadRequestResponse()
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(strings.ReplaceAll(conserv.Split(err.Field()), " ", "_"))
		var value interface{}
		if val, ok := err.Value().(string); ok && val == "" {
			value = "null or empty"
		} else {
			value = err.Value()
		}
		errResponse.Details = append(errResponse.Details,
			fmt.Sprintf("invalid field '%s' with value %v", field, value))
	}

	ctx.JSON(http.StatusBadRequest, errResponse)
}

func setInternalErrorResponse(ctx *gin.Context, err error) {
	errResponse := errorResponse{
		Name:    "internal_error",
		Message: "internal server error",
		Details: []string{err.Error()},
	}

	ctx.JSON(http.StatusInternalServerError, errResponse)
}

type server struct {
	createProjectUsecase  create.ImportProjectUsecase
	analyzeProjectUsecase analyze.AnalyzeProjectUsecase
}

func (s server) createAnalysis(ctx *gin.Context) {

}
