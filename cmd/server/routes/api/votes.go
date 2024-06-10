package api

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/vote"

	"github.com/gin-gonic/gin"
)

func InitializeVotes(router *gin.Engine) {
	router.GET("/api/votes", middleware.IsAuthenticated, vote.GetVotes)
	router.POST("/api/votes/:projectID", middleware.IsAuthenticated,vote.CreateVote)
	router.GET("/api/votes/:id", middleware.IsAuthenticated,vote.GetVoteByID)
	router.PUT("/api/votes/:id",  middleware.IsAuthenticated,vote.UpdateVote)
	router.DELETE("/api/votes/:id", middleware.IsAuthenticated,vote.DeleteVote)
}