package ui

import (
	"go-api/cmd/server/middleware"
	"net/http"
	"time"

	cmt "go-api/internal/comment"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func CommentForm(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.Redirect(http.StatusFound, "/auth/login")
		}
		targetID := c.Param("targetID")
		_, err := primitive.ObjectIDFromHex(targetID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
			return
		}
		c.HTML(200, "form-create-comment.html", gin.H{ 
			"session_user": user,
			"targetID": targetID,
		})
}

func CommentSubmitForm(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	var errors = []gin.H{}

	var comment = cmt.Comment{}
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	targetID := c.Param("targetID")
	target, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.HTML(200, "response-create-comment.html", gin.H{
			"submitted": false,
			"message": "Error creating comment",
			"errors": errors,
		})
	}

	comment.TargetID = target
	comment.AuthorID = user.(map[string]interface{})["nickname"].(string)
	comment.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	
	if !cmt.RequiredFields(comment) {
		c.HTML(200, "response-create-comment.html", gin.H{
			"submitted": false,
			"message": "Error creating comment",
			"errors": errors,
		})
		return
	}
	

	id, err := cmt.DbCreateComment(comment)
	if err != nil {
		c.HTML(200, "response-create-comment.html", gin.H{
			"submitted": false,
			"message": "Error creating comment",
			"errors": errors,
		})
		return
	}

	comment.ID = id
	c.HTML(200, "response-create-comment.html", gin.H{
		"submitted": true,
		"message": "Comment created successfully",
		"errors": errors,
	})
}

func InitializeCommentsUI(router *gin.Engine) {
	router.GET("/ui/comment/form/:targetID", middleware.IsAuthenticated, CommentForm)
	router.POST("/ui/comment/submit/:targetID", middleware.IsAuthenticated, CommentSubmitForm)
}