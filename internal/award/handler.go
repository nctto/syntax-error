package award

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAwards(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := AwardsDefaultQueryParams(c)
	awards, err := DbGetAllAwards(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, awards)
}

func GetAwardByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	award, err := DbGetAwardID(id)
	if err != nil {
		fmt.Println("Error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": "Award not found"})
		return
	}
	c.JSON(http.StatusOK, award)
}

func CreateAward(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	projectID := c.Param("projectID")
	id, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid project ID"})
		return
	}

	awardTypeID := c.Param("typeID")
	typeID, err := primitive.ObjectIDFromHex(awardTypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid type ID"})
		return
	}

	_, err = DbAwardTypeExists(typeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var award = Award{}
	award.TargetID = id
	award.TypeID = typeID
	award.AuthorID = user.(map[string]interface{})["id"].(string)
	award.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = DbCreateAward(award)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	award.ID = id
	c.JSON(http.StatusCreated, award)
}

func DeleteAward(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteAward(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Award deleted"})
}