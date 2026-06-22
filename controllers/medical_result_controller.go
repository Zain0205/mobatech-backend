package controllers

import (
	"backend/models"
	"backend/services"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results})
}

// GET /api/medical-results
func (c *MedicalResultController) GetUserResults(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat.(float64))

	results, err := c.service.GetUserMedicalResults(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": results})
}

// GET /api/medical-results/:id
func (c *MedicalResultController) GetByID(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	result, err := c.service.GetMedicalResultByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// POST /api/admin/medical-results
func (c *MedicalResultController) Create(ctx *gin.Context) {
	var req models.MedicalResult
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.service.CreateMedicalResult(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Medical result created successfully", "data": result})
}

// PUT /api/admin/medical-results/:id
func (c *MedicalResultController) Update(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	
	var req models.MedicalResult
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = uint(id)

	result, err := c.service.UpdateMedicalResult(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Medical result updated successfully", "data": result})
}

// DELETE /api/admin/medical-results/:id
func (c *MedicalResultController) Delete(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err := c.service.DeleteMedicalResult(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Medical result deleted successfully"})
}
