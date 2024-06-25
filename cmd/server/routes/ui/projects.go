package ui

import (
	"go-api/cmd/server/middleware"
	pr "go-api/internal/project"

	"github.com/gin-gonic/gin"
)

func InitializeProjectsUI(router *gin.Engine) {
	router.GET("/ui/projects/all", pr.GetHTMLAllProjects)
	router.GET("/ui/projects/form", middleware.IsAuthenticated, pr.GetHTMLSubmitProjectForm)
	router.POST("/ui/projects/form/submit", middleware.IsAuthenticated, pr.GetHTMLSubmitProjectForm)
}