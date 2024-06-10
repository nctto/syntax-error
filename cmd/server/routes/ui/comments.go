package ui

import (
	"go-api/cmd/server/middleware"
	"net/http"

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

func InitializeCommentsUI(router *gin.Engine) {
	router.GET("/ui/comment/form/:targetID", middleware.IsAuthenticated, CommentForm)
	router.POST("/ui/comment/submit/:targetID", middleware.IsAuthenticated, cmt.CreateComment)
}