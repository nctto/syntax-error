package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pr "go-api/internal/project"

	"github.com/gin-contrib/sessions"
)

func InitializeHomePage(router *gin.Engine) {
	router.GET("/", func (c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := pr.ProjectsDefaultQueryParams(c)
	projects, err := pr.DbGetAllProjects(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	projectsView := pr.ProjectsToProjectView(projects)
	c.HTML(200, "home-page.html", gin.H{
		"title": "syntax error", 
		"session_user": user,
		"projects": projectsView,
		"page": page,
		"limit": limit,
		"nextPage": "2",
		"sortBy": sortBy,
	})
})
}