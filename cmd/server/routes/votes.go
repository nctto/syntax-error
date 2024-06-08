package routes

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/vote"

	"github.com/gin-gonic/gin"
)

func InitializeVotes(router *gin.Engine) {
	router.GET("/votes", middleware.IsAuthenticated, vote.GetVotes)
	router.POST("/votes/:projectID", middleware.IsAuthenticated,vote.CreateVote)
	router.GET("/votes/:id", middleware.IsAuthenticated,vote.GetVoteByID)
	router.PUT("/votes/:id",  middleware.IsAuthenticated,vote.UpdateVote)
	router.DELETE("/votes/:id", middleware.IsAuthenticated,vote.DeleteVote)
}