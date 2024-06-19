package ui

import (
	"net/http"

	"go-api/cmd/server/middleware"
	pr "go-api/internal/project"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
)

func UiGetAllProjects(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := pr.ProjectsDefaultQueryParams(c)
	projects, total,  err := pr.DbGetAllProjects(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	paginatedProjects := pr.ProjectPaginatedView(projects, total, page, limit, sortBy)

	c.HTML(200, "list-projects.html", gin.H{
		"session_user": user,
		"projects": paginatedProjects,
	})
}

func UiSubmitProjectForm(c *gin.Context) {
		
	session := sessions.Default(c)
	user := session.Get("profile")

	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}

	var errors = []gin.H{}
	var newProject pr.Project
	if err := c.BindJSON(&newProject); err != nil {
		errors = append(errors, gin.H{"message": err.Error()})
	}
	newProject.AuthorID = user.(map[string]interface{})["nickname"].(string)
	if !pr.RequiredFields(newProject) {
		errors = append(errors, gin.H{"message": "Missing required fields"})
	}

	if len(errors) > 0 {
		c.HTML(200, "response-create-project.html", gin.H{
			"submitted": false,
			"message": "Error creating project",
			"errors": errors,
		})
		return
	}

	id, err := pr.DbCreateProject(newProject)
	if err != nil {
		errors = append(errors, gin.H{"message": err.Error()})
	}
	newProject.ID = id
	c.HTML(200, "response-create-project.html", gin.H{
		"submitted": true,
		"message": "Project created successfully",
		"errors": errors,
	})
}

func InitializeProjectsUI(router *gin.Engine) {
	router.GET("/ui/projects/all", UiGetAllProjects)
	router.GET("/ui/projects/form", middleware.IsAuthenticated, UiSubmitProjectForm)
	router.POST("/ui/projects/form/submit", middleware.IsAuthenticated, UiSubmitProjectForm)
}