package main

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database Connection
	config.ConnectDatabase()

	// Run Auto Migrations for our example User model
	config.DB.AutoMigrate(&models.User{})

	// Initialize Gin Router
	r := gin.Default()

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

	// Start the server
	r.Run(":8080")
}
