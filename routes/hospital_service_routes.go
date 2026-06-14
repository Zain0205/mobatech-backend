package routes

import (
	"backend/controllers"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupHospitalServiceRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repositories.NewHospitalServiceRepository(db)
	service := services.NewHospitalServiceService(repo)
	controller := controllers.NewHospitalServiceController(service)

	// Public API for mobile
	publicGroup := r.Group("/api/services")
	{
		publicGroup.GET("", controller.GetAll)
		publicGroup.GET("/:id", controller.GetByID)
	}

	// Admin API
	adminGroup := r.Group("/api/admin/services")
	{
		adminGroup.POST("", controller.Create)
		adminGroup.PUT("/:id", controller.Update)
		adminGroup.DELETE("/:id", controller.Delete)
	}
}
