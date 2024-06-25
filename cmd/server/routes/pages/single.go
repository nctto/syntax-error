package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	cm "go-api/internal/comment"
	pr "go-api/internal/project"

	"github.com/gin-contrib/sessions"
)

func InitializeSingleProjectPage(router *gin.Engine) {
	router.GET("/:targetID", func (c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")
		
		targetID := c.Param("targetID")
		id, err := primitive.ObjectIDFromHex(targetID)
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

		projectView := pr.ProjectToProjectView(project)
		c.HTML(200, "single-project-page.html", gin.H{
			"title": "syntax error", 
			"ID": targetID,
			"TargetID": targetID,
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
			"Comments": comments,
		})
})
}