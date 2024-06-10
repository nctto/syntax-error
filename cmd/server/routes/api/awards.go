package api

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/award"

	"github.com/gin-gonic/gin"
)

func InitializeAwards(router *gin.Engine) {
	router.GET("/api/awards/:projectID", middleware.IsAuthenticated, award.GetAwards)
	router.POST("/api/awards/:projectID/:typeID", middleware.IsAuthenticated, award.CreateAward)
	router.GET("/api/awards/:id", middleware.IsAuthenticated, award.GetAwardByID)
	router.DELETE("/api/awards/:id", middleware.IsAuthenticated, award.DeleteAward)
}