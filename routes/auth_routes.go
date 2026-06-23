package routes

import (
	"backend/controllers"
	"backend/repositories"
	"backend/services"

	"backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repositories.NewAuthRepository(db)
	service := services.NewAuthService(repo)
	authController := controllers.NewAuthController(service)
	profileController := controllers.NewProfileController(service)

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
		authGroup.GET("/me", middleware.AuthMiddleware(), authController.Me)
	}

	userGroup := r.Group("/api/users")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.PUT("/profile", profileController.UpdateProfile)
		userGroup.POST("/family-members", profileController.AddFamilyMember)
		userGroup.DELETE("/family-members/:id", profileController.DeleteFamilyMember)
	}
}
