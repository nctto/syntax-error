package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-api/cmd/server/middleware"
	cm "go-api/internal/comment"
	pr "go-api/internal/project"

	"github.com/gin-contrib/sessions"
)

const (
	title = "syntax error"
)

func Home(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := pr.ProjectsDefaultQueryParams(c)
	projects, err := pr.DbGetAllProjects(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	projectsView := pr.ProjectsToProjectView(projects)
	c.HTML(200, "home.html", gin.H{
		"title": title, 
		"session_user": user,
		"projects": projectsView,
		"page": page,
		"limit": limit,
		"nextPage": "2",
		"sortBy": sortBy,
	})
}


func CreateForm(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.Redirect(http.StatusFound, "/auth/login")
		}
		c.HTML(200, "create.html", gin.H{
			"title": title, 
			"session_user": user,
		})
}

func SingleProject(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.Redirect(http.StatusFound, "/auth/login")
		}

		projectID := c.Param("projectID")
		id, err := primitive.ObjectIDFromHex(projectID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
			return
		}

		project, err := pr.DbGetProjectID(id, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		comments, err := cm.DbGetAllComments(1, 10, "best", user, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		commentsView := cm.CommentsToCommentView(comments)

		projectView := pr.ProjectToProjectView(project)
		c.HTML(200, "single.html", gin.H{
			"title": title, 
			"session_user": user,
			"Title": projectView.Title,
			"Content": projectView.Content,
			"AuthorID": projectView.AuthorID,
			"CreatedAt": projectView.CreatedAt,
			"VotesTotal": projectView.VotesTotal,
			"Voted": projectView.Voted,
			"CommentsTotal": projectView.CommentsTotal,
			"Awards": projectView.Awards,
			"AwardsTotal": projectView.AwardsTotal,
			"Tags": projectView.Tags,
			"Comments": commentsView,
		})
}


func InitializeHome(router *gin.Engine) {
	router.GET("/", Home)
	router.GET("/create", middleware.IsAuthenticated, CreateForm)
	router.GET("/:projectID", SingleProject)
}