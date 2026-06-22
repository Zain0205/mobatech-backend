package controllers

import (
	"backend/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ForYouController struct {
	service services.ForYouService
}

func NewForYouController(service services.ForYouService) *ForYouController {
	return &ForYouController{service}
}

func (ctrl *ForYouController) GetRecommendations(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHENTICATED"})
		return
	}

	userID := fmt.Sprintf("%v", userIDStr)
	articles, err := ctrl.service.GenerateRecommendations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": articles})
}
