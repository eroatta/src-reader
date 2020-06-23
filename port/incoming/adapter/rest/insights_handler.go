package rest

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/eroatta/src-reader/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type createInsightsCommand struct {
	AnalysisID string `json:"analysis_id" validate:"uuid"`
}

type insightsResponse struct {
	AnalysisID string                  `json:"analysis_id"`
	ProjectRef string                  `json:"project_ref"`
	Summary    insightsSummaryResponse `json:"identifiers"`
	Overall    float64                 `json:"accuracy"`
	Packages   []packageResponse       `json:"packages"`
}

type insightsSummaryResponse struct {
	Total    int `json:"total"`
	Exported int `json:"exported"`
}

type packageResponse struct {
	Name    string                  `json:"name"`
	Summary insightsSummaryResponse `json:"identifiers"`
	Ratio   float64                 `json:"accuracy"`
	Files   []string                `json:"files"`
}

// RegisterGainInsightsUsecase sets the endpoint and the handler on the REST service to
// handle the insights creation.
func RegisterGainInsightsUsecase(r *gin.Engine, uc usecase.GainInsightsUsecase) *gin.Engine {
	r.POST("/insights", func(c *gin.Context) {
		createInsights(c, uc)
	})
	// r.GET("/insigths/$id", handler)
	// r.DELETE("/insights/$id", handler)

	return r
}

func createInsights(ctx *gin.Context, uc usecase.GainInsightsUsecase) {
	var cmd createInsightsCommand

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

	insights, err := uc.Process(ctx, uuid.MustParse(cmd.AnalysisID))
	switch err {
	case nil:
		// do nothing
	case usecase.ErrPreviousInsightsFound:
		setBadRequestResponse(ctx, fmt.Errorf("previous insights exist for analysis with ID: %v", cmd.AnalysisID))
		return
	case usecase.ErrIdentifiersNotFound:
		setBadRequestResponse(ctx, fmt.Errorf("non-existing identifiers for analysis ID %v", cmd.AnalysisID))
		return
	default:
		setInternalErrorResponse(ctx, err)
		return
	}

	response := insightsResponse{
		AnalysisID: cmd.AnalysisID,
		Summary:    insightsSummaryResponse{},
		Packages:   make([]packageResponse, 0),
	}
	weighted := 0.0

	for _, insight := range insights {
		response.ProjectRef = insight.ProjectRef
		response.Summary.Total += insight.TotalIdentifiers
		response.Summary.Exported += insight.TotalExported
		weighted += insight.TotalWeight

		files := make([]string, 0)
		for file := range insight.Files {
			files = append(files, file)
		}
		sort.Strings(files)

		response.Packages = append(response.Packages, packageResponse{
			Name: insight.Package,
			Summary: insightsSummaryResponse{
				Total:    insight.TotalIdentifiers,
				Exported: insight.TotalExported,
			},
			Ratio: insight.Rate(),
			Files: files,
		})
	}
	response.Overall = weighted / float64(response.Summary.Total)

	ctx.JSON(http.StatusCreated, response)
}
