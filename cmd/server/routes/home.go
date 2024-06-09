package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api/cmd/server/middleware"
	"go-api/internal/project"

	"github.com/gin-contrib/sessions"
)

const (
	title = "syntax error"
)

func Home(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := project.ProjectsDefaultQueryParams(c)
	projects, err := project.DbGetAllProjects(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	projectsView := project.ProjectsToProjectView(projects)
	c.HTML(200, "index.html", gin.H{
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
			"project": project.FakeProject(),
		})
}


func InitializeHome(router *gin.Engine) {
	router.GET("/", Home)
	router.GET("/create", middleware.IsAuthenticated, CreateForm)
}