package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MedicalResultController struct {
	service services.MedicalResultService
}

func NewMedicalResultController(service services.MedicalResultService) *MedicalResultController {
	return &MedicalResultController{service}
}

// GET /api/admin/medical-results
func (c *MedicalResultController) GetAll(ctx *gin.Context) {
	results, err := c.service.GetAllMedicalResults()
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", results))
}

// GET /api/medical-results
func (c *MedicalResultController) GetUserResults(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}
	userID := uint(userIDFloat.(float64))

	results, err := c.service.GetUserMedicalResults(userID)
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", results))
}

// GET /api/medical-results/:id
func (c *MedicalResultController) GetByID(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	result, err := c.service.GetMedicalResultByID(uint(id))
	if err != nil {
		ctx.Error(utils.NewAppError(utils.ErrNotFound, http.StatusNotFound, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", result))
}

// POST /api/admin/medical-results
func (c *MedicalResultController) Create(ctx *gin.Context) {
	var req models.MedicalResult
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	result, err := c.service.CreateMedicalResult(&req)
	if err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", result))
}

// PUT /api/admin/medical-results/:id
func (c *MedicalResultController) Update(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	var req models.MedicalResult
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}
	req.ID = uint(id)

	result, err := c.service.UpdateMedicalResult(&req)
	if err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", result))
}

// DELETE /api/admin/medical-results/:id
func (c *MedicalResultController) Delete(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err := c.service.DeleteMedicalResult(uint(id)); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}
