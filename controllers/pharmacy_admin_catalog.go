package controllers

import (
	"backend/models"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c *PharmacyController) AdminCreateCategory(ctx *gin.Context) {
	var cat models.MedicineCategory
	if err := ctx.ShouldBindJSON(&cat); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}
	if err := c.service.CreateCategory(&cat); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, utils.BuildSuccess("CREATED", "Resource created successfully", cat))
}

func (c *PharmacyController) AdminUpdateCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var cat models.MedicineCategory
	if err := ctx.ShouldBindJSON(&cat); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}
	cat.ID = uint(id)

	if err := c.service.UpdateCategory(&cat); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", cat))
}

func (c *PharmacyController) AdminDeleteCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := c.service.DeleteCategory(uint(id)); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}

func (c *PharmacyController) AdminCreateMedicine(ctx *gin.Context) {
	var med models.Medicine
	if err := ctx.ShouldBindJSON(&med); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}
	if err := c.service.CreateMedicine(&med); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, utils.BuildSuccess("CREATED", "Resource created successfully", med))
}

func (c *PharmacyController) AdminUpdateMedicine(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var med models.Medicine
	if err := ctx.ShouldBindJSON(&med); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}
	med.ID = uint(id)

	if err := c.service.UpdateMedicine(&med); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", med))
}

func (c *PharmacyController) AdminDeleteMedicine(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := c.service.DeleteMedicine(uint(id)); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}
