package routes

import (
	"backend/controllers"
	"backend/middleware"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPatientSupportRoutes(r *gin.Engine, db *gorm.DB) {
	// Medical Results DI
	mrRepo := repositories.NewMedicalResultRepository(db)
	mrService := services.NewMedicalResultService(mrRepo)
	mrController := controllers.NewMedicalResultController(mrService)

	// Reminder DI
	remRepo := repositories.NewReminderRepository(db)
	remService := services.NewReminderService(remRepo)
	remController := controllers.NewReminderController(remService)

	// Medical Results Routes
	mrUserGroup := r.Group("/api/medical-results")
	mrUserGroup.Use(middleware.AuthMiddleware())
	{
		mrUserGroup.GET("", mrController.GetUserResults)
		mrUserGroup.GET("/:id", mrController.GetByID)
	}

	mrAdminGroup := r.Group("/api/admin/medical-results")
	{
		mrAdminGroup.GET("", mrController.GetAll)
		mrAdminGroup.POST("", mrController.Create)
		mrAdminGroup.PUT("/:id", mrController.Update)
		mrAdminGroup.DELETE("/:id", mrController.Delete)
	}

	// Reminders Routes
	remUserGroup := r.Group("/api/reminders")
	remUserGroup.Use(middleware.AuthMiddleware())
	{
		remUserGroup.GET("", remController.GetUserReminders)
		remUserGroup.GET("/unread-count", remController.GetUnreadCount)
		remUserGroup.PUT("/:id/read", remController.MarkAsRead)
	}

	remAdminGroup := r.Group("/api/admin/reminders")
	{
		remAdminGroup.GET("", remController.GetAll)
		remAdminGroup.POST("", remController.Create)
		remAdminGroup.DELETE("/:id", remController.Delete)
	}
}
