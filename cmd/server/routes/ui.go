package routes

import (
	"net/http"
	"time"

	"go-api/cmd/server/middleware"
	pr "go-api/internal/project"
	vt "go-api/internal/vote"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-contrib/sessions"
)

func UiGetProjects(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := pr.ProjectsDefaultQueryParams(c)
	projects, err := pr.DbGetAllProjects(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	projectsView := pr.ProjectsToProjectView(projects)

	c.HTML(200, "projects.html", gin.H{
		"title": title, 
		"session_user": user,
		"projects": projectsView,
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
	var newProject pr.Project
	if err := c.BindJSON(&newProject); err != nil {
		errors = append(errors, gin.H{"message": err.Error()})
	}

	if !pr.RequiredFields(newProject) {
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

	id, err := pr.DbCreateProject(newProject)
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

func UiCreateVote(c *gin.Context) {
	
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	
	projectId := c.Param("projectID")
	id, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid vote ID"})
		return
	}
	
	exists ,err := pr.DbProjectExists(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Project not found"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Project not found"})
		return
	}

	
	var newVote vt.Vote
	newVote.ProjectID = id
	newVote.AuthorID = user.(map[string]interface{})["nickname"].(string)
	newVote.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	vote,err := vt.DbVoteExists(newVote.ProjectID, newVote.AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	
	if vote {
		err = vt.DbDeleteVoteByAuthor(newVote.ProjectID, newVote.AuthorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		numberOfVotes, err := vt.DbGetProjectVotes(newVote.ProjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.HTML(200, "vote.html", gin.H{
			"ID": projectId, 
			"Votes": numberOfVotes,
			"Voted": false,
		})
		return
	}

	if _,err := vt.DbCreateVote(newVote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	numberOfVotes, err := vt.DbGetProjectVotes(newVote.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.HTML(201, "vote.html", gin.H{
		"ID": projectId, 
		"Votes": numberOfVotes,
		"Voted": true,
	})
}

func InitializeUI(router *gin.Engine) {
	router.GET("/ui/projects", UiGetProjects)
	router.POST("/ui/projects", middleware.IsAuthenticated, UiCreateProject)
	router.POST("/ui/votes/:projectID", middleware.IsAuthenticated, UiCreateVote)
}