package controllers

import (
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c *PharmacyController) GetMyPrescriptions(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}
	userID := uint(userIDVal.(float64))

	prescriptions, err := c.service.GetPrescriptionsByUserID(userID)
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", prescriptions))
}
func (c *PharmacyController) GetPrescriptionDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	order, err := c.authorizeOrderAccess(ctx, id)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", order))
}
func (c *PharmacyController) CreateOrder(ctx *gin.Context) {
	var req models.PharmacyOrder
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}
	req.UserID = uint(userIDVal.(float64))

	if err := c.service.CreateOrder(&req); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.BuildSuccess("CREATED", "Resource created successfully", req))
}
func (c *PharmacyController) GetMyOrders(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}
	userID := uint(userIDVal.(float64))

	orders, err := c.service.GetOrdersByUserID(userID)
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", orders))
}
func (c *PharmacyController) GetOrderDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	order, err := c.authorizeOrderAccess(ctx, id)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", order))
}
func (c *PharmacyController) CancelOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	order, err := c.authorizeOrderAccess(ctx, id)
	if err != nil {
		return
	}

	if order.Status != "Pending" {
		ctx.Error(utils.NewValidationError("Only pending orders can be cancelled"))
		return
	}

	if err := c.service.UpdateOrderStatus(uint(id), "Cancelled"); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}

func (c *PharmacyController) authorizeOrderAccess(ctx *gin.Context, id int) (*models.PharmacyOrder, error) {
	order, err := c.service.GetOrderByID(uint(id))
	if err != nil {
		ctx.Error(utils.NewAppError(utils.ErrNotFound, http.StatusNotFound, "Order not found"))
		return nil, err
	}

	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		err = utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized")
		ctx.Error(err)
		return nil, err
	}
	userID := uint(userIDVal.(float64))
	if order.UserID != userID {
		err = fmt.Errorf("Access denied")
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return nil, err
	}
	return order, nil
}
