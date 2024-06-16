package api

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/vote"

	"github.com/gin-gonic/gin"
)

func InitializeVotes(router *gin.Engine) {
	router.GET("/api/votes", middleware.IsAuthenticated, vote.GetVotes)
	router.POST("/api/vote/:targetID", middleware.IsAuthenticated,vote.CreateVote)
}