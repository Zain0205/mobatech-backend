package routes

import (
	"backend/controllers"
	"backend/middleware"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupChatRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repositories.NewChatRepository(db)
	service := services.NewChatService(repo)
	controller := controllers.NewChatController(service)

	chatGroup := r.Group("/api/chat")
	chatGroup.Use(middleware.AuthMiddleware())
	{
		chatGroup.GET("/sessions", controller.GetUserSessions)
		chatGroup.POST("/sessions", controller.CreateSession)
		chatGroup.GET("/sessions/:id/messages", controller.GetSessionMessages)
		chatGroup.POST("/sessions/:id/stream", controller.StreamChat)
	}
}
