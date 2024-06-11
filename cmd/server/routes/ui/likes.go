package ui

import (
	"net/http"
	"time"

	"go-api/cmd/server/middleware"
	lk "go-api/internal/like"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-contrib/sessions"
)

func UiSubmitLike(c *gin.Context) {
	
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	
	targetID := c.Param("targetID")
	id, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid like ID"})
		return
	}
	
	var newLike lk.Like
	newLike.TargetID = id
	newLike.AuthorID = user.(map[string]interface{})["nickname"].(string)
	newLike.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	like,err := lk.DbLikeExists(newLike.TargetID, newLike.AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	
	if like {
		err = lk.DbDeleteLikeByAuthor(newLike.TargetID, newLike.AuthorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.HTML(200, "component-like.html", gin.H{
			"ID": targetID, 
			"Liked": false,
		})
		return
	}

	if _,err := lk.DbCreateLike(newLike); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.HTML(201, "component-like.html", gin.H{
		"ID": targetID, 
		"Liked": true,
	})
}

func InitializeLikesUI(router *gin.Engine) {
	router.POST("/ui/likes/submit/:targetID", middleware.IsAuthenticated, UiSubmitLike)
}