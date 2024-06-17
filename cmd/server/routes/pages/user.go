package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api/internal/project"
	usr "go-api/internal/user"

	"github.com/gin-contrib/sessions"
)

func InitializeSingleUserPage(router *gin.Engine) {
	router.GET("/u/:username", func (c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		
		username := c.Param("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username is required"})
			return
		}

		profile, err := usr.DbGetUserByUsername(username)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		projects, err := project.DbGetProjectsByUser(profile.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		profileView := usr.UserToUserView(profile)
		c.HTML(200, "single-user-page.html", gin.H{
			"session_user": user,
			"title": profileView.Username, 
			"Username": profileView.Username,
			"CreatedAt": profileView.CreatedAt,
			"projects": project.ProjectsToProjectView(projects),
		})
})
}