package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/cmd/server/authenticator"

	"go-api/internal/project"

	"github.com/gin-contrib/sessions"
)

const (
	title = "syntax error"
)

func InitializeHome(router *gin.Engine, a *authenticator.Authenticator) {
	
	router.GET("/", func(c *gin.Context) {

		session := sessions.Default(c)
		user := session.Get("profile")

		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")
		sortBy := c.DefaultQuery("sort_by", "top")

		projects, err := project.DbGetAllProjects(page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.HTML(200, "index.html", gin.H{
			"title": title, 
			"session_user": user,
			"projects": projects,
			"page": page,
			"limit": limit,
			"nextPage": "2",
			"sortBy": sortBy,
		})
	})
	
	router.GET("/ui/projects", func(c *gin.Context) {

		session := sessions.Default(c)
		user := session.Get("profile")

		page := c.Query("page")
		limit := c.Query("limit")
		sortBy := c.Query("sort_by")

		fmt.Println("Page",page)
		fmt.Println("Limit",limit)

		projects, err := project.DbGetAllProjects(page, limit)
		fmt.Println("All Projects",len(projects))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		nextPage, _ := strconv.Atoi(page)
		nextPage++


		c.HTML(200, "projects.html", gin.H{
			"title": title, 
			"session_user": user,
			"projects": projects,
			"page": page,
			"limit": limit,
			"nextPage": nextPage,
			"sortBy": sortBy,
		})
	})


	router.GET("/create", func(c *gin.Context) {
		
		session := sessions.Default(c)
		user := session.Get("profile")

		if user == nil {
			c.Redirect(http.StatusFound, "/auth/login")
		}

		c.HTML(200, "create.html", gin.H{
			"title": "syntax error", 
			"session_user": user,
			"project": project.FakeProject(),
		})
	})


	router.POST("/ui/create/", func(c *gin.Context) {
		
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
				"title": "syntax error", 
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
			"title": "syntax error", 
			"submitted": true,
			"message": "Project created successfully",
			"errors": errors,
		})

	})
}