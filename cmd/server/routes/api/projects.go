package api

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/project"

	"github.com/gin-gonic/gin"
)

func InitializeProjects(router *gin.Engine) {
	router.GET("/api/projects", middleware.IsAuthenticated, project.GetProjects)
	router.POST("/api/projects", middleware.IsAuthenticated, project.CreateProject)
	router.GET("/api/projects/:id", middleware.IsAuthenticated, project.GetProjectByID)
	router.PUT("/api/projects/:id",  middleware.IsAuthenticated, project.UpdateProject)
	router.DELETE("/api/projects/:id", middleware.IsAuthenticated, project.DeleteProject)
	// router.GET("/api/projects/fake",middleware.IsAuthenticated, project.FakeProjects)
	// router.GET("/api/projects/random", middleware.IsAuthenticated, project.GetRandomProject)
}