package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EmergencyController struct {
	service services.EmergencyService
}

func NewEmergencyController(service services.EmergencyService) *EmergencyController {
	return &EmergencyController{service}
}

func (c *EmergencyController) SubmitRequest(ctx *gin.Context) {
	var req models.EmergencyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	req.UserID = uint(userID.(float64))

	if err := c.service.CreateRequest(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-dispatch after 3 seconds in a background goroutine
	go func(emergencyID uint) {
		time.Sleep(3 * time.Second)
		_ = c.service.UpdateStatus(emergencyID, "Dispatched")
	}(req.ID)

	ctx.JSON(http.StatusCreated, req)
}

func (c *EmergencyController) GetUserHistory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	history, err := c.service.GetHistoryByUser(uint(userID.(float64)))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, history)
}

func (c *EmergencyController) GetAllAdmin(ctx *gin.Context) {
	reqs, err := c.service.GetAllRequests()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reqs)
}

func (c *EmergencyController) UpdateStatusAdmin(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateStatus(uint(id), req.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
