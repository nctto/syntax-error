package like

import (
	"go-api/internal/comment"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetLikes(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")

	likes, err := DbGetAllLikes(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, likes)
}

func GetLikeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	like, err := DbGetLikeID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Like not found"})
		return
	}
	c.JSON(http.StatusOK, like)
}

func CreateLike(c *gin.Context) {
	
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	
	commentId := c.Param("commentID")
	id, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid like ID"})
		return
	}
	
	comment,err := comment.DbGetCommentID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Comment not found"})
		return
	}
	
	var newLike Like
	newLike.TargetID = comment.ID
	newLike.AuthorID = user.(map[string]interface{})["nickname"].(string)
	newLike.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	like,err := DbLikeExists(newLike.TargetID, newLike.AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	
	if like {
		err = DbDeleteLikeByAuthor(newLike.TargetID, newLike.AuthorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		numberOfLikes, err := DbGetCommentLikes(newLike.TargetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ID": commentId ,"Likes": numberOfLikes, "Liked": false})
		return
	}

	if _,err := DbCreateLike(newLike); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	numberOfLikes, err := DbGetCommentLikes(newLike.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"ID": commentId ,"Likes": numberOfLikes, "Liked": true})
}

func UpdateLike(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	var updatedLike Like
	if err := c.BindJSON(&updatedLike); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = DbUpdateLike(id, updatedLike)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	updatedLike.ID = id
	c.JSON(http.StatusOK, updatedLike)
}

func DeleteLike(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteLike(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Like deleted"})
}
