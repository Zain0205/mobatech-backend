package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AppointmentController struct {
	appointmentService services.AppointmentService
}

func NewAppointmentController(appointmentService services.AppointmentService) *AppointmentController {
	return &AppointmentController{appointmentService}
}

// GET /api/admin/appointments
func (c *AppointmentController) GetAllAppointments(ctx *gin.Context) {
	appointments, err := c.appointmentService.GetAllAppointments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": appointments})
}

// GET /api/appointments
func (c *AppointmentController) GetUserAppointments(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat.(float64))

	appointments, err := c.appointmentService.GetUserAppointments(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": appointments})
}

// POST /api/appointments
func (c *AppointmentController) BookAppointment(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat.(float64))

	var req models.Appointment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := c.appointmentService.BookAppointment(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Appointment booked successfully", "data": appointment})
}

// POST /api/appointments/:id/cancel
func (c *AppointmentController) CancelAppointment(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat.(float64))

	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	
	// Assuming non-admin endpoint for users to cancel their own appointments
	if err := c.appointmentService.CancelAppointment(uint(id), userID, false); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled successfully"})
}

// POST /api/admin/appointments/:id/cancel
func (c *AppointmentController) AdminCancelAppointment(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	
	// Admin can cancel any appointment
	if err := c.appointmentService.CancelAppointment(uint(id), 0, true); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled successfully by admin"})
}

// POST /api/admin/appointments/:id/approve
func (c *AppointmentController) ApproveAppointment(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err := c.appointmentService.ApproveAppointment(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Appointment approved successfully"})
}
