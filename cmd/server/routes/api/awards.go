package api

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/award"

	"github.com/gin-gonic/gin"
)

func InitializeAwards(router *gin.RouterGroup) {
	router.GET("/:id", middleware.IsAuthenticated, award.GetAwardByID)
	router.POST("/:id//:postID", middleware.IsAuthenticated, award.CreateAward)
	router.DELETE("/awards/:id", middleware.IsAuthenticated, award.DeleteAward)
}