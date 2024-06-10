package api

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/comment"

	"github.com/gin-gonic/gin"
)

func InitializeComments(router *gin.Engine) {
	router.POST("/api/comments/:targetID", middleware.IsAuthenticated, comment.CreateComment)
	router.DELETE("/api/comments/:id", middleware.IsAuthenticated, comment.DeleteComment)
}