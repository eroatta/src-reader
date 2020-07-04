package rest

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
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

// RegisterCreateProjectUsecase defines the proper URI and HTTP method to execute the CreateProjectUsecase.
func RegisterCreateProjectUsecase(r *gin.Engine, uc usecase.CreateProjectUsecase) *gin.Engine {
	r.POST("/projects", func(c *gin.Context) {
		createProject(c, uc)
	})

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

	project, err := uc.Process(ctx, cmd.Reference)
	if err != nil {
		log.WithError(err).Error("unexpected error executing createProjectUsecase")
		setInternalErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toProjectResponse(project))
}

// RegisterGetProjectUsecase defines the proper URI and HTTP method to execute the GetProjectUsecase.
func RegisterGetProjectUsecase(r *gin.Engine, uc usecase.GetProjectUsecase) *gin.Engine {
	r.GET("/projects/:id", func(c *gin.Context) {
		getProject(c, uc)
	})

	return r
}

func getProject(ctx *gin.Context, uc usecase.GetProjectUsecase) {
	ID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		setNotFoundResponse(ctx, fmt.Errorf("project with ID: %s can't be found", ctx.Param("id")))
		return
	}

	project, err := uc.Process(ctx, ID)
	switch err {
	case nil:
		// do nothing
	case usecase.ErrProjectNotFound:
		setNotFoundResponse(ctx, fmt.Errorf("project with ID: %s can't be found", ID.String()))
		return
	default:
		log.WithError(err).Error("unexpected error executing getProjectUsecase")
		setInternalErrorResponse(ctx, fmt.Errorf("error accessing project with ID: %s", ID.String()))
		return
	}

	ctx.JSON(http.StatusOK, toProjectResponse(project))
}

// RegisterDeleteProjectUsecase defines the proper URI and HTTP method to execute the DeleteProjectUsecase.
func RegisterDeleteProjectUsecase(r *gin.Engine, uc usecase.DeleteProjectUsecase) *gin.Engine {
	r.DELETE("/projects/:id", func(c *gin.Context) {
		deleteProject(c, uc)
	})

	return r
}

func deleteProject(ctx *gin.Context, uc usecase.DeleteProjectUsecase) {
	projectID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		setNotFoundResponse(ctx, fmt.Errorf("project %s can't be found", ctx.Param("id")))
		return
	}

	err = uc.Process(ctx, projectID)
	if err != nil && err != usecase.ErrProjectNotFound {
		log.WithError(err).Error("unexpected error executing deleteProjectUsecase")
		setInternalErrorResponse(ctx, fmt.Errorf("error deleting project with ID: %v", projectID))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func toProjectResponse(project entity.Project) projectResponse {
	return projectResponse{
		ID:        project.ID.String(),
		Status:    project.Status,
		Reference: project.Reference,
		Metadata: metadataResponse{
			RemoteID:      project.Metadata.RemoteID,
			Owner:         project.Metadata.Owner,
			Fullname:      project.Metadata.Fullname,
			Description:   project.Metadata.Description,
			CloneURL:      project.Metadata.CloneURL,
			DefaultBranch: project.Metadata.DefaultBranch,
			License:       project.Metadata.License,
			CreatedAt:     project.Metadata.CreatedAt,
			UpdatedAt:     project.Metadata.UpdatedAt,
			IsFork:        project.Metadata.IsFork,
			Size:          project.Metadata.Size,
			Stargazers:    project.Metadata.Stargazers,
			Watchers:      project.Metadata.Watchers,
			Forks:         project.Metadata.Forks,
		},
		SourceCode: sourcecodeResponse{
			Hash:     project.SourceCode.Hash,
			Location: project.SourceCode.Location,
			Files:    project.SourceCode.Files,
		},
	}
}
