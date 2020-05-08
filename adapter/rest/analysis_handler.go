package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eroatta/src-reader/usecase/analyze"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type createAnalysisCommand struct {
	Repository string `json:"repository" validate:"url"`
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

func (s server) createAnalysis(ctx *gin.Context) {
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

	analysis, err := s.analyzeProjectUsecase.Analyze(ctx, cmd.Repository)
	switch err {
	case nil:
		// do nothing
	case analyze.ErrProjectNotFound:
		setBadRequestResponse(ctx, fmt.Errorf("non-existing repository %s", cmd.Repository))
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
