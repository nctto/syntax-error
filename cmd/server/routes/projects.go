package routes

import (
	"go-api/internal/project"

	"go-api/cmd/server/authenticator"

	"github.com/gin-gonic/gin"
)

func InitializeProjects(router *gin.Engine, a *authenticator.Authenticator) {
	router.GET("/projects", project.GetProjects)
	router.GET("/projects/fake", project.FakeProjects)
	router.POST("/projects", project.CreateProject)
	router.GET("/projects/random", project.GetRandomProject)
	router.GET("/projects/:id", project.GetProjectByID)
	router.PUT("/projects/:id",  project.UpdateProject)
	router.DELETE("/projects/:id", project.DeleteProject)
}