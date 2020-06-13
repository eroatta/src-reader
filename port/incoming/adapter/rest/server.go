package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eroatta/token/conserv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// requestValidator represents a validator capable of analyzing the values of the incoming
// request bodies.
var requestValidator = validator.New()

func NewServer() *gin.Engine {
	r := gin.Default()

	setMetricsCollectors(r)

	r.GET("/ping", pingHandler)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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

func setBadRequestResponse(ctx *gin.Context, err error) {
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

func setNotFoundResponse(ctx *gin.Context, err error) {
	errResponse := errorResponse{
		Name:    "not_found",
		Message: "resource not found",
		Details: make([]string, 0),
	}
	errResponse.Details = append(errResponse.Details, err.Error())

	ctx.JSON(http.StatusNotFound, errResponse)
}

func setInternalErrorResponse(ctx *gin.Context, err error) {
	errResponse := errorResponse{
		Name:    "internal_error",
		Message: "internal server error",
		Details: []string{err.Error()},
	}

	ctx.JSON(http.StatusInternalServerError, errResponse)
}
