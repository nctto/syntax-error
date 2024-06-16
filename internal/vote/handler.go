package vote

import (
	"fmt"
	"go-api/internal/comment"
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

func SubmitVote(target string, targetID primitive.ObjectID, authorID string) (int32, bool, error) {

	if target == "p" {
		exists, err := project.DbProjectExists(targetID)
		if err != nil {
			return 0, false, err
		}
		if !exists {
			return 0, false, fmt.Errorf("project not found")
		}	
	}

	if target == "c" {
		exists, err := comment.DbCommentExists(targetID)
		if err != nil {
			return 0, false, err
		}
		if !exists {
			return 0, false, fmt.Errorf("comment not found")
		}
	
	}

	fmt.Print("TargetID: ", targetID)
	var voted bool
	var v Vote
	v.TargetID = targetID
	v.AuthorID = authorID
	v.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	vote, err := DbVoteExists(v.TargetID, v.AuthorID)
	if err != nil {
		return 0, false, err
	}
	fmt.Println("Vote exists: ", vote)
	if vote {
		err = DbDeleteVoteByAuthor(v.TargetID, v.AuthorID)
		if err != nil {
			return 0, true, err
		}
		voted = false
	} else { 
		_,err := DbCreateVote(v)
		if err != nil {
			return 0, false, err
		}
		voted = true
	}
	votes, err := DbGetTargetVotes(v.TargetID)
	if err != nil {
		return 0, voted, err
	}
	return votes, voted, nil
}

func CreateVote(c *gin.Context) {
	
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
	
	votes, voted, err := SubmitVote("p", id, user.(map[string]interface{})["nickname"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ID": projectId ,"Votes": votes, "Voted": voted})
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
