package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
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
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", appointments))
}

// GET /api/appointments
func (c *AppointmentController) GetUserAppointments(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}
	userID := uint(userIDFloat.(float64))

	appointments, err := c.appointmentService.GetUserAppointments(userID)
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", appointments))
}

// POST /api/appointments
func (c *AppointmentController) BookAppointment(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}
	userID := uint(userIDFloat.(float64))

	var req models.Appointment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	appointment, err := c.appointmentService.BookAppointment(userID, &req)
	if err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", appointment))
}

// POST /api/appointments/:id/cancel
func (c *AppointmentController) CancelAppointment(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}
	userID := uint(userIDFloat.(float64))

	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	// Assuming non-admin endpoint for users to cancel their own appointments
	if err := c.appointmentService.CancelAppointment(uint(id), userID, false); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}

// POST /api/admin/appointments/:id/cancel
func (c *AppointmentController) AdminCancelAppointment(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	// Admin can cancel any appointment
	if err := c.appointmentService.CancelAppointment(uint(id), 0, true); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}

// POST /api/admin/appointments/:id/approve
func (c *AppointmentController) ApproveAppointment(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err := c.appointmentService.ApproveAppointment(uint(id)); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}
