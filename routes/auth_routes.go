package routes

import (
	"backend/controllers"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repositories.NewAuthRepository(db)
	service := services.NewAuthService(repo)
	controller := controllers.NewAuthController(service)

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", controller.Register)
		authGroup.POST("/login", controller.Login)
	}
}
