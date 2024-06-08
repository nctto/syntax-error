package routes

import (
	"net/http"

	"go-api/cmd/server/middleware"
	"go-api/internal/project"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
)

func UiGetProjects(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := project.ProjectsDefaultQueryParams(c)
	projects, err := project.DbGetAllProjects(page, limit, sortBy, user)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.HTML(200, "projects.html", gin.H{
		"title": title, 
		"session_user": user,
		"projects": projects,
		"page": page,
		"limit": limit,
		"nextPage": page + 1,
		"sortBy": sortBy,
	})
}

func UiCreateProject(c *gin.Context) {
		
	session := sessions.Default(c)
	user := session.Get("profile")

	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}

	var errors = []gin.H{}
	var newProject project.Project
	if err := c.BindJSON(&newProject); err != nil {
		errors = append(errors, gin.H{"message": err.Error()})
	}

	if !project.RequiredFields(newProject) {
		errors = append(errors, gin.H{"message": "Missing required fields"})
	}

	if len(errors) > 0 {
		c.HTML(200, "submitted.html", gin.H{
			"title": title, 
			"submitted": false,
			"message": "Error creating project",
			"errors": errors,
		})
		return
	}

	id, err := project.DbCreateProject(newProject)
	if err != nil {
		errors = append(errors, gin.H{"message": err.Error()})
	}
	newProject.ID = id
	c.HTML(200, "submitted.html", gin.H{
		"title": title, 
		"submitted": true,
		"message": "Project created successfully",
		"errors": errors,
	})
}

func InitializeUI(router *gin.Engine) {
	router.GET("/ui/projects", UiGetProjects)
	router.POST("/ui/create/", middleware.IsAuthenticated, UiCreateProject)
}