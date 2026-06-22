package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReminderController struct {
	service services.ReminderService
}

func NewReminderController(service services.ReminderService) *ReminderController {
	return &ReminderController{service}
}

// GET /api/admin/reminders
func (c *ReminderController) GetAll(ctx *gin.Context) {
	reminders, err := c.service.GetAllReminders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": reminders})
}

// GET /api/reminders
func (c *ReminderController) GetUserReminders(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat.(float64))

	reminders, err := c.service.GetUserReminders(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": reminders})
}

// GET /api/reminders/unread-count
func (c *ReminderController) GetUnreadCount(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat.(float64))

	count, err := c.service.GetUnreadCount(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{"count": count}})
}

// POST /api/admin/reminders
func (c *ReminderController) Create(ctx *gin.Context) {
	var req models.Reminder
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminder, err := c.service.CreateReminder(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reminder created successfully", "data": reminder})
}

// PUT /api/reminders/:id/read
func (c *ReminderController) MarkAsRead(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat.(float64))

	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	
	if err := c.service.MarkAsRead(uint(id), userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reminder marked as read successfully"})
}

// DELETE /api/admin/reminders/:id
func (c *ReminderController) Delete(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err := c.service.DeleteReminder(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reminder deleted successfully"})
}
