package routes

import (
	"backend/controllers"
	"backend/middleware"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupForYouRoutes(r *gin.Engine, db *gorm.DB) {
	chatRepo := repositories.NewChatRepository(db)
	service := services.NewForYouService(chatRepo)
	controller := controllers.NewForYouController(service)

	group := r.Group("/api/for-you")
	group.Use(middleware.AuthMiddleware())
	{
		group.GET("/recommendations", controller.GetRecommendations)
	}
}
