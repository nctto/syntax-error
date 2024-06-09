package vote

import (
	"go-api/internal/project"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetVotes(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")

	votes, err := DbGetAllVotes(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, votes)
}

func GetVoteByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	vote, err := DbGetVoteID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Vote not found"})
		return
	}
	c.JSON(http.StatusOK, vote)
}

func CreateVote(c *gin.Context) {
	
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
	
	project,err := project.DbGetProjectID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Project not found"})
		return
	}
	
	var newVote Vote
	newVote.ProjectID = project.ID
	newVote.AuthorID = user.(map[string]interface{})["nickname"].(string)
	newVote.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	vote,err := DbVoteExists(newVote.ProjectID, newVote.AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	
	if vote {
		err = DbDeleteVoteByAuthor(newVote.ProjectID, newVote.AuthorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		numberOfVotes, err := DbGetProjectVotes(newVote.ProjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ID": projectId ,"Votes": numberOfVotes, "Voted": false})
		return
	}

	if _,err := DbCreateVote(newVote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	numberOfVotes, err := DbGetProjectVotes(newVote.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"ID": projectId ,"Votes": numberOfVotes, "Voted": true})
}

func UpdateVote(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	var updatedVote Vote
	if err := c.BindJSON(&updatedVote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = DbUpdateVote(id, updatedVote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	updatedVote.ID = id
	c.JSON(http.StatusOK, updatedVote)
}

func DeleteVote(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteVote(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vote deleted"})
}
