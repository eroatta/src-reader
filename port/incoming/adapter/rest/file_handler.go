package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eroatta/src-reader/usecase/file"
	"github.com/gin-gonic/gin"
)

func RegisterOriginalFileUsecase(r *gin.Engine, uc file.OriginalFileUsecase) *gin.Engine {
	r.GET("/files/originals/:owner/:project/*file", func(c *gin.Context) {
		getOriginalFile(c, uc)
	})

	return r
}

func getOriginalFile(ctx *gin.Context, uc file.OriginalFileUsecase) {
	projectRef := fmt.Sprintf("%s/%s", ctx.Param("owner"), ctx.Param("project"))
	raw, err := uc.Process(ctx, projectRef, strings.TrimPrefix(ctx.Param("file"), "/"))
	switch err {
	case nil:
		// do nothing
	case file.ErrProjectNotFound:
		ctx.AbortWithStatus(http.StatusNotFound)
	case file.ErrFileNotFound:
		ctx.AbortWithStatus(http.StatusNotFound)
	default:
		ctx.String(http.StatusInternalServerError, "Ooops! Something went wrong...")
		return
	}

	ctx.String(200, string(raw))
}

func RegisterRewrittenFileUsecase(r *gin.Engine, uc file.RewrittenFileUsecase) *gin.Engine {
	r.GET("/files/rewritten/:owner/:project/*file", func(c *gin.Context) {
		getRewrittenFile(c, uc)
	})

	return r
}

func getRewrittenFile(ctx *gin.Context, uc file.OriginalFileUsecase) {
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
