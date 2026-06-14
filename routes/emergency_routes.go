package routes

import (
	"backend/controllers"
	"backend/middleware"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupEmergencyRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repositories.NewEmergencyRepository(db)
	service := services.NewEmergencyService(repo)
	controller := controllers.NewEmergencyController(service)
	trackingController := controllers.NewTrackingController(service)

	// WebSocket tracking route (no auth - WebSocket clients have trouble with auth headers)
	r.GET("/api/emergencies/:id/track", trackingController.TrackAmbulance)

	// User Routes (Protected)
	userGroup := r.Group("/api/emergencies")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.POST("", controller.SubmitRequest)
		userGroup.GET("/history", controller.GetUserHistory)
	}

	// Admin Routes (Can be protected by AdminMiddleware if you have one, using basic for now)
	adminGroup := r.Group("/api/admin/emergencies")
	{
		adminGroup.GET("", controller.GetAllAdmin)
		adminGroup.PUT("/:id/status", controller.UpdateStatusAdmin)
	}
}
