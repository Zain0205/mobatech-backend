package routes

import (
	"fmt"
	"backend/controllers"
	"backend/middleware"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupDoctorRoutes(router *gin.Engine, db *gorm.DB) {
	fmt.Println("SetupDoctorRoutes called")
	// Repositories
	doctorRepo := repositories.NewDoctorRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	appointmentRepo := repositories.NewAppointmentRepository(db)

	// Services
	doctorService := services.NewDoctorService(doctorRepo)
	scheduleService := services.NewScheduleService(scheduleRepo)
	appointmentService := services.NewAppointmentService(appointmentRepo, scheduleRepo)

	// Controllers
	doctorController := controllers.NewDoctorController(doctorService)
	scheduleController := controllers.NewScheduleController(scheduleService)
	appointmentController := controllers.NewAppointmentController(appointmentService)

	api := router.Group("/api")

	// Mobile Endpoints
	api.GET("/doctors", doctorController.GetDoctors)
	api.GET("/doctors/:id", doctorController.GetDoctorByID)
	api.GET("/doctors/:id/schedules", scheduleController.GetDoctorSchedules)

	// Protected Mobile Endpoints (require authentication)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/appointments", appointmentController.GetUserAppointments)
		protected.POST("/appointments", appointmentController.BookAppointment)
		protected.POST("/appointments/:id/cancel", appointmentController.CancelAppointment)
	}

	// Admin Endpoints
	admin := api.Group("/admin")
	// If you have admin middleware, add it here. E.g. admin.Use(middleware.AdminMiddleware())
	// For now using AuthMiddleware or just public if no specific admin auth is provided
	// In production, should verify admin role
	// admin.Use(middleware.AuthMiddleware()) 
	{
		// Doctor CRUD
		admin.POST("/doctors", doctorController.CreateDoctor)
		admin.PUT("/doctors/:id", doctorController.UpdateDoctor)
		admin.DELETE("/doctors/:id", doctorController.DeleteDoctor)

		// Schedule CRUD
		admin.POST("/schedules", scheduleController.CreateSchedule)
		admin.PUT("/schedules/:id", scheduleController.UpdateSchedule)
		admin.DELETE("/schedules/:id", scheduleController.DeleteSchedule)

		// Appointment Management
		admin.GET("/appointments", appointmentController.GetAllAppointments)
		admin.POST("/appointments/:id/approve", appointmentController.ApproveAppointment)
		admin.POST("/appointments/:id/cancel", appointmentController.AdminCancelAppointment)
	}
}
