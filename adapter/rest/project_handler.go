package rest

import "github.com/gin-gonic/gin"

type postCreateProjectCommand struct {
	Repository string `json:"repository"`
}

func (s server) createProject(ctx *gin.Context) {

}
