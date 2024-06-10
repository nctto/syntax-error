package ui

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

func UiSubmitVote(c *gin.Context) {
	
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
		c.HTML(200, "component-vote.html", gin.H{
			"ID": projectId, 
			"VotesTotal": numberOfVotes,
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
	c.HTML(201, "component-vote.html", gin.H{
		"ID": projectId, 
		"VotesTotal": numberOfVotes,
		"Voted": true,
	})
}

func InitializeVotesUI(router *gin.Engine) {
	router.POST("/ui/votes/submit/:projectID", middleware.IsAuthenticated, UiSubmitVote)
}