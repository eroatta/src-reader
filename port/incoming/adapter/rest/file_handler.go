package rest

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/eroatta/src-reader/usecase/file"
	"github.com/gin-gonic/gin"
)

func RegisterOriginalFileUsecase(r *gin.Engine, uc file.OriginalFileUsecase) *gin.Engine {
	r.GET("/files/originals/:owner/:project/*file", func(c *gin.Context) {
		getFile(c, uc)
	})

	return r
}

func RegisterRewrittenFileUsecase(r *gin.Engine, uc file.RewrittenFileUsecase) *gin.Engine {
	r.GET("/files/rewritten/:owner/:project/*file", func(c *gin.Context) {
		getFile(c, uc)
	})

	return r
}

type fileUsecase interface {
	Process(ctx context.Context, projectRef string, filename string) ([]byte, error)
}

func getFile(ctx *gin.Context, uc fileUsecase) {
	projectRef := fmt.Sprintf("%s/%s", ctx.Param("owner"), ctx.Param("project"))
	raw, err := uc.Process(ctx, projectRef, strings.TrimPrefix(ctx.Param("file"), "/"))
	switch err {
	case nil:
		// do nothing
	case file.ErrProjectNotFound:
		ctx.AbortWithStatus(http.StatusNotFound)
	case file.ErrFileNotFound:
		ctx.AbortWithStatus(http.StatusNotFound)
	case file.ErrIdentifiersNotFound:
		ctx.String(http.StatusConflict, "Ooops! Did you already analyze this project?")
		return
	default:
		ctx.String(http.StatusInternalServerError, "Ooops! Something went wrong...")
		return
	}

	ctx.String(200, string(raw))
}
