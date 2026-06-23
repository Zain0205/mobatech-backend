package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PolyclinicController struct {
	service services.PolyclinicService
}

func NewPolyclinicController(service services.PolyclinicService) *PolyclinicController {
	return &PolyclinicController{service}
}

// GET /api/polyclinics
func (c *PolyclinicController) GetPolyclinics(ctx *gin.Context) {
	polys, err := c.service.GetAllPolyclinics()
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", polys))
}

// GET /api/polyclinics/:id
func (c *PolyclinicController) GetPolyclinicByID(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	poly, err := c.service.GetPolyclinicByID(uint(id))
	if err != nil {
		ctx.Error(utils.NewAppError(utils.ErrNotFound, http.StatusNotFound, "Polyclinic not found"))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", poly))
}

// POST /api/admin/polyclinics
func (c *PolyclinicController) CreatePolyclinic(ctx *gin.Context) {
	var req models.Polyclinic
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}
	if err := c.service.CreatePolyclinic(&req); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", req))
}

// PUT /api/admin/polyclinics/:id
func (c *PolyclinicController) UpdatePolyclinic(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	var req models.Polyclinic
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}
	req.ID = uint(id)
	if err := c.service.UpdatePolyclinic(&req); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", req))
}

// DELETE /api/admin/polyclinics/:id
func (c *PolyclinicController) DeletePolyclinic(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err := c.service.DeletePolyclinic(uint(id)); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}

// POST /api/admin/polyclinics/:id/schedules
func (c *PolyclinicController) CreateSchedule(ctx *gin.Context) {
	polyID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	var req models.PolyclinicSchedule
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}
	req.PolyclinicID = uint(polyID)
	if err := c.service.CreateSchedule(&req); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", req))
}

// PUT /api/admin/polyclinics/schedules/:sched_id
func (c *PolyclinicController) UpdateSchedule(ctx *gin.Context) {
	schedID, _ := strconv.ParseUint(ctx.Param("sched_id"), 10, 32)
	var req models.PolyclinicSchedule
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}
	req.ID = uint(schedID)
	if err := c.service.UpdateSchedule(&req); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", req))
}

// DELETE /api/admin/polyclinics/schedules/:sched_id
func (c *PolyclinicController) DeleteSchedule(ctx *gin.Context) {
	schedID, _ := strconv.ParseUint(ctx.Param("sched_id"), 10, 32)
	if err := c.service.DeleteSchedule(uint(schedID)); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}
