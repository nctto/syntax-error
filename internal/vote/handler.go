package vote

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetVotes(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")

	projects, err := DbGetAllVotes(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func GetVoteByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	project, err := DbGetVoteID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Vote not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func CreateVote(c *gin.Context) {

	var newVote Vote
	if err := c.BindJSON(&newVote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !RequiredFields(newVote) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields"})
		return
	}

	id, err := DbCreateVote(newVote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newVote.ID = id
	c.JSON(http.StatusCreated, newVote)
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
