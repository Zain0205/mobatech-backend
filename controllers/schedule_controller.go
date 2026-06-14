package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ScheduleController struct {
	scheduleService services.ScheduleService
}

func NewScheduleController(scheduleService services.ScheduleService) *ScheduleController {
	return &ScheduleController{scheduleService}
}

// GET /api/doctors/:id/schedules
func (c *ScheduleController) GetDoctorSchedules(ctx *gin.Context) {
	doctorID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	schedules, err := c.scheduleService.GetDoctorSchedules(uint(doctorID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": schedules})
}

// POST /api/admin/schedules
func (c *ScheduleController) CreateSchedule(ctx *gin.Context) {
	var input models.DoctorSchedule
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.scheduleService.CreateSchedule(&input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Schedule created successfully", "data": input})
}

// PUT /api/admin/schedules/:id
func (c *ScheduleController) UpdateSchedule(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	var input models.DoctorSchedule
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := c.scheduleService.UpdateSchedule(uint(id), &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Schedule updated successfully", "data": schedule})
}

// DELETE /api/admin/schedules/:id
func (c *ScheduleController) DeleteSchedule(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err := c.scheduleService.DeleteSchedule(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
}
