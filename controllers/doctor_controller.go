package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DoctorController struct {
	doctorService services.DoctorService
}

func NewDoctorController(doctorService services.DoctorService) *DoctorController {
	return &DoctorController{doctorService}
}

// GET /api/doctors?specialization=...
func (c *DoctorController) GetDoctors(ctx *gin.Context) {
	specialization := ctx.Query("specialization")
	doctors, err := c.doctorService.GetAllDoctors(specialization)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": doctors})
}

// GET /api/doctors/:id
func (c *DoctorController) GetDoctorByID(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	doctor, err := c.doctorService.GetDoctorByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": doctor})
}

// POST /api/admin/doctors
func (c *DoctorController) CreateDoctor(ctx *gin.Context) {
	var doctor models.Doctor
	if err := ctx.ShouldBindJSON(&doctor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.doctorService.CreateDoctor(&doctor); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Doctor created successfully", "data": doctor})
}

// PUT /api/admin/doctors/:id
func (c *DoctorController) UpdateDoctor(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	var input models.Doctor
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor, err := c.doctorService.UpdateDoctor(uint(id), &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Doctor updated successfully", "data": doctor})
}

// DELETE /api/admin/doctors/:id
func (c *DoctorController) DeleteDoctor(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err := c.doctorService.DeleteDoctor(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Doctor deleted successfully"})
}
