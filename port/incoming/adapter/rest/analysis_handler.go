package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eroatta/src-reader/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type createAnalysisCommand struct {
	ProjectID string `json:"project_id" validate:"uuid"`
}

type analysisResponse struct {
	ID          string          `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	ProjectRef  string          `json:"project_ref"`
	Miners      []string        `json:"miners"`
	Splitters   []string        `json:"splitters"`
	Expanders   []string        `json:"expanders"`
	Files       summaryResponse `json:"files_summary"`
	Identifiers summaryResponse `json:"identifiers_summary"`
}

type summaryResponse struct {
	Total        int      `json:"total"`
	Valid        int      `json:"valid"`
	Failed       int      `json:"failed"`
	ErrorSamples []string `json:"error_samples"`
}

func RegisterAnalyzeProjectUsecase(r *gin.Engine, uc usecase.AnalyzeProjectUsecase) *gin.Engine {
	r.POST("/analysis", func(c *gin.Context) {
		createAnalysis(c, uc)
	})
	// r.GET("/analysis/$id", internal.getAnalysis)
	// r.DELETE("/analysis/$id", internal.deleteAnalysis)
	// r.GET("/analysis/$id/identifiers", internal.getIdentifiers)

	return r
}

func createAnalysis(ctx *gin.Context, uc usecase.AnalyzeProjectUsecase) {
	var cmd createAnalysisCommand

	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		log.WithError(err).Debug("failed to bind JSON body")
		setBadRequestResponse(ctx, err)
		return
	}

	if err := requestValidator.Struct(cmd); err != nil {
		log.WithError(err).Debug("failed while validating the command")
		setBadRequestOnValidationResponse(ctx, err)
		return
	}

	analysis, err := uc.Process(ctx, uuid.MustParse(cmd.ProjectID))
	switch err {
	case nil:
		// do nothing
	case usecase.ErrProjectNotFound:
		setBadRequestResponse(ctx, fmt.Errorf("project with ID: %s can't be found", cmd.ProjectID))
		return
	default:
		setInternalErrorResponse(ctx, err)
		return
	}

	response := analysisResponse{
		ID:         analysis.ID,
		ProjectRef: analysis.ProjectName,
		CreatedAt:  analysis.DateCreated,
		Miners:     analysis.PipelineMiners,
		Splitters:  analysis.PipelineSplitters,
		Expanders:  analysis.PipelineExpanders,
		Files: summaryResponse{
			Total:        analysis.FilesTotal,
			Valid:        analysis.FilesValid,
			Failed:       analysis.FilesError,
			ErrorSamples: analysis.FilesErrorSamples,
		},
		Identifiers: summaryResponse{
			Total:        analysis.IdentifiersTotal,
			Valid:        analysis.IdentifiersValid,
			Failed:       analysis.IdentifiersError,
			ErrorSamples: analysis.IdentifiersErrorSamples,
		},
	}
	ctx.JSON(http.StatusCreated, response)
}
