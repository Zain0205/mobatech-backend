package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HospitalServiceController struct {
	service services.HospitalServiceService
}

func NewHospitalServiceController(service services.HospitalServiceService) *HospitalServiceController {
	return &HospitalServiceController{service}
}

func (c *HospitalServiceController) GetAll(ctx *gin.Context) {
	services, err := c.service.GetAll()
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", services))
}

func (c *HospitalServiceController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	service, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.Error(utils.NewAppError(utils.ErrNotFound, http.StatusNotFound, "Service not found"))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", service))
}

func (c *HospitalServiceController) Create(ctx *gin.Context) {
	var req models.HospitalService
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	if err := c.service.Create(&req); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, utils.BuildSuccess("CREATED", "Resource created successfully", req))
}

func (c *HospitalServiceController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	service, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.Error(utils.NewAppError(utils.ErrNotFound, http.StatusNotFound, "Service not found"))
		return
	}

	if err := ctx.ShouldBindJSON(&service); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	if err := c.service.Update(service); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", service))
}

func (c *HospitalServiceController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}
