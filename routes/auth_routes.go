package routes

import (
	"backend/controllers"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"backend/middleware"
)

func SetupAuthRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repositories.NewAuthRepository(db)
	service := services.NewAuthService(repo)
	controller := controllers.NewAuthController(service)

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", controller.Register)
		authGroup.POST("/login", controller.Login)
		authGroup.GET("/me", middleware.AuthMiddleware(), controller.Me)
	}

	userGroup := r.Group("/api/users")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.PUT("/profile", controller.UpdateProfile)
		userGroup.POST("/family-members", controller.AddFamilyMember)
		userGroup.DELETE("/family-members/:id", controller.DeleteFamilyMember)
	}
}
