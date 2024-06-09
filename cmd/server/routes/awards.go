package routes

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/award"

	"github.com/gin-gonic/gin"
)

func InitializeAwards(router *gin.Engine) {
	router.GET("/awards/:projectID", middleware.IsAuthenticated, award.GetAwards)
	router.POST("/awards/:projectID/:typeID", middleware.IsAuthenticated, award.CreateAward)
	router.GET("/awards/:id", middleware.IsAuthenticated, award.GetAwardByID)
	router.DELETE("/awards/:id", middleware.IsAuthenticated, award.DeleteAward)
}