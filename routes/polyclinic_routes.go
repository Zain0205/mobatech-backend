package routes

import (
	"backend/controllers"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPolyclinicRoutes(router *gin.Engine, db *gorm.DB) {
	repo := repositories.NewPolyclinicRepository(db)
	service := services.NewPolyclinicService(repo)
	controller := controllers.NewPolyclinicController(service)

	api := router.Group("/api")

	// Mobile public endpoints
	api.GET("/polyclinics", controller.GetPolyclinics)
	api.GET("/polyclinics/:id", controller.GetPolyclinicByID)

	// Admin endpoints
	admin := api.Group("/admin")
	{
		admin.POST("/polyclinics", controller.CreatePolyclinic)
		admin.PUT("/polyclinics/:id", controller.UpdatePolyclinic)
		admin.DELETE("/polyclinics/:id", controller.DeletePolyclinic)

		admin.POST("/polyclinics/:id/schedules", controller.CreateSchedule)
		admin.PUT("/polyclinics/schedules/:sched_id", controller.UpdateSchedule)
		admin.DELETE("/polyclinics/schedules/:sched_id", controller.DeleteSchedule)
	}
}
