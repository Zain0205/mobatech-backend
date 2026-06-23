package controllers

import (
	"backend/models"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c *PharmacyController) AdminCreatePrescription(ctx *gin.Context) {
	var p models.Prescription
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}
	if err := c.service.CreatePrescription(&p); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, utils.BuildSuccess("CREATED", "Resource created successfully", p))
}

func (c *PharmacyController) AdminGetAllPrescriptions(ctx *gin.Context) {
	prescriptions, err := c.service.GetAllPrescriptions()
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", prescriptions))
}

func (c *PharmacyController) AdminGetAllOrders(ctx *gin.Context) {
	orders, err := c.service.GetAllOrders()
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", orders))
}

func (c *PharmacyController) AdminUpdateOrderStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	if err := c.service.UpdateOrderStatus(uint(id), req.Status); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}

func (c *PharmacyController) AdminUpdateOrderPayment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req struct {
		PaymentStatus string `json:"payment_status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	if err := c.service.UpdateOrderPayment(uint(id), req.PaymentStatus); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}
