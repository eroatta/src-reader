package rest

import (
	"net/http"
	"regexp"
	"time"

	"github.com/eroatta/src-reader/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
)

func init() {
	regex := regexp.MustCompile(`^[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+$`)
	err := requestValidator.RegisterValidation("reference", func(fl validator.FieldLevel) bool {
		return regex.MatchString(fl.Field().String())
	})
	if err != nil {
		log.WithError(err).Panic("unable to configure validators")
	}
}

type postCreateProjectCommand struct {
	Reference string `json:"reference" validate:"reference"`
}

type projectResponse struct {
	ID         string             `json:"id"`
	Status     string             `json:"status"`
	Reference  string             `json:"reference"`
	Metadata   metadataResponse   `json:"metadata"`
	SourceCode sourcecodeResponse `json:"source_code"`
}

type metadataResponse struct {
	RemoteID      string     `json:"remote_id"`
	Owner         string     `json:"owner"`
	Fullname      string     `json:"fullname"`
	Description   string     `json:"description"`
	CloneURL      string     `json:"clone_url"`
	DefaultBranch string     `json:"branch"`
	License       string     `json:"license"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	IsFork        bool       `json:"is_fork"`
	Size          int32      `json:"size"`
	Stargazers    int32      `json:"stargazers"`
	Watchers      int32      `json:"watchers"`
	Forks         int32      `json:"forks"`
}

type sourcecodeResponse struct {
	Hash     string   `json:"hash"`
	Location string   `json:"location"`
	Files    []string `json:"files"`
}

func RegisterCreateProjectUsecase(r *gin.Engine, uc usecase.CreateProjectUsecase) *gin.Engine {
	r.POST("/projects", func(c *gin.Context) {
		createProject(c, uc)
	})
	// r.GET("/projects/$id", internal.getProject)
	// r.DELETE("/projects/$id", internal.deleteProject)

	return r
}

func createProject(ctx *gin.Context, uc usecase.CreateProjectUsecase) {
	var cmd postCreateProjectCommand

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

	p, err := uc.Process(ctx, cmd.Reference)
	if err != nil {
		log.WithError(err).Error("unexpected error executing createProjectUsecase")
		setInternalErrorResponse(ctx, err)
		return
	}

	response := projectResponse{
		ID:        p.ID,
		Status:    p.Status,
		Reference: p.Reference,
		Metadata: metadataResponse{
			RemoteID:      p.Metadata.RemoteID,
			Owner:         p.Metadata.Owner,
			Fullname:      p.Metadata.Fullname,
			Description:   p.Metadata.Description,
			CloneURL:      p.Metadata.CloneURL,
			DefaultBranch: p.Metadata.DefaultBranch,
			License:       p.Metadata.License,
			CreatedAt:     p.Metadata.CreatedAt,
			UpdatedAt:     p.Metadata.UpdatedAt,
			IsFork:        p.Metadata.IsFork,
			Size:          p.Metadata.Size,
			Stargazers:    p.Metadata.Stargazers,
			Watchers:      p.Metadata.Watchers,
			Forks:         p.Metadata.Forks,
		},
		SourceCode: sourcecodeResponse{
			Hash:     p.SourceCode.Hash,
			Location: p.SourceCode.Location,
			Files:    p.SourceCode.Files,
		},
	}
	ctx.JSON(http.StatusCreated, response)
}
