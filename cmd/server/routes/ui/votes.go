package ui

import (
	"net/http"

	"go-api/cmd/server/middleware"
	vt "go-api/internal/vote"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-contrib/sessions"
)

func UiSubmitVote(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	projectId := c.Param("targetID")
	id, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid vote ID"})
		return
	}

	nickname := user.(map[string]interface{})["nickname"].(string)
	votes, voted, err := vt.SubmitVote("project", id, nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.HTML(201, "component-vote.html", gin.H{
		"ID":         projectId,
		"VotesTotal": votes,
		"Voted":      voted,
	})
}

func InitializeVotesUI(router *gin.Engine) {
	router.POST("/ui/vote/submit/:targetID", middleware.IsAuthenticated, UiSubmitVote)
}

