package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
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
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", schedules))
}

// POST /api/admin/schedules
func (c *ScheduleController) CreateSchedule(ctx *gin.Context) {
	var input models.DoctorSchedule
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	if err := c.scheduleService.CreateSchedule(&input); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", input))
}

// PUT /api/admin/schedules/:id
func (c *ScheduleController) UpdateSchedule(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	var input models.DoctorSchedule
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	schedule, err := c.scheduleService.UpdateSchedule(uint(id), &input)
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", schedule))
}

// DELETE /api/admin/schedules/:id
func (c *ScheduleController) DeleteSchedule(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err := c.scheduleService.DeleteSchedule(uint(id)); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}
