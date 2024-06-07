package routes

import (
	"go-api/internal/vote"

	"go-api/cmd/server/authenticator"

	"github.com/gin-gonic/gin"
)

func InitializeVotes(router *gin.Engine, a *authenticator.Authenticator) {
	router.GET("/votes", vote.GetVotes)
	router.POST("/votes", vote.CreateVote)
	router.GET("/votes/:id", vote.GetVoteByID)
	router.PUT("/votes/:id",  vote.UpdateVote)
	router.DELETE("/votes/:id", vote.DeleteVote)
}