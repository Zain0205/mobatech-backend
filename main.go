package main

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"backend/routes"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database Connection
	config.ConnectDatabase()

	// Run Auto Migrations for our example User model and Chat models
	config.DB.AutoMigrate(&models.User{}, &models.ChatSession{}, &models.ChatMessage{})

	// Initialize Gin Router
	r := gin.Default()
	r.Use(middleware.ErrorHandler())
	r.Use(cors.Default())

	r.Static("/uploads", "./uploads")

	// Health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Sample endpoint to get users
	r.GET("/users", func(c *gin.Context) {
		var users []models.User
		config.DB.Find(&users)
		c.JSON(http.StatusOK, gin.H{"data": users})
	})

	// Setup Auth Routes
	routes.SetupAuthRoutes(r, config.DB)
	// Setup Chat Routes
	routes.SetupChatRoutes(r, config.DB)
	// Setup Hospital Service Routes
	routes.SetupHospitalServiceRoutes(r, config.DB)
	// Setup Emergency Routes
	routes.SetupEmergencyRoutes(r, config.DB)
	// Setup Pharmacy Routes
	routes.SetupPharmacyRoutes(r, config.DB)
	// Setup Doctor Routes
	routes.SetupDoctorRoutes(r, config.DB)
	// Setup Polyclinic Routes
	routes.SetupPolyclinicRoutes(r, config.DB)
	// Setup Patient Support Routes
	routes.SetupPatientSupportRoutes(r, config.DB)
	// Setup For You Routes
	routes.SetupForYouRoutes(r, config.DB)

	// Start the server
	r.Run(":8080")
}
