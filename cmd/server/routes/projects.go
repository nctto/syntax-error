package routes

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/project"

	"github.com/gin-gonic/gin"
)

func InitializeProjects(router *gin.Engine) {
	router.GET("/projects", middleware.IsAuthenticated, project.GetProjects)
	router.GET("/projects/fake",middleware.IsAuthenticated, project.FakeProjects)
	router.POST("/projects", middleware.IsAuthenticated, project.CreateProject)
	router.GET("/projects/random", middleware.IsAuthenticated, project.GetRandomProject)
	router.GET("/projects/:id", middleware.IsAuthenticated, project.GetProjectByID)
	router.PUT("/projects/:id",  middleware.IsAuthenticated, project.UpdateProject)
	router.DELETE("/projects/:id", middleware.IsAuthenticated, project.DeleteProject)
}