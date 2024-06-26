package comment

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetComments(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	targetID := c.Param("targetID")
	id, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid target ID"})
		return
	}

	page, limit, sortBy := CommentsDefaultQueryParams(c)
	comments, err := DbGetAllComments(page, limit, sortBy, user, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)

}

func CreateComment(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var comment = Comment{}
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	targetID := c.Param("targetID")
	target, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid target ID"})
		return
	}

	comment.TargetID = target
	comment.AuthorID = user.(map[string]interface{})["nickname"].(string)
	comment.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	
	if !RequiredFields(comment) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields"})
		return
	}
	

	id, err := DbCreateComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	comment.ID = id
	c.JSON(http.StatusCreated, comment)
}

func DeleteComment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteComment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}