package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	service services.AuthService
}

func NewProfileController(service services.AuthService) *ProfileController {
	return &ProfileController{service}
}

func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}

	fullName := ctx.PostForm("full_name")
	phone := ctx.PostForm("phone_number")
	bloodType := ctx.PostForm("blood_type")
	height, _ := strconv.Atoi(ctx.PostForm("height"))
	weight, _ := strconv.Atoi(ctx.PostForm("weight"))
	allergies := ctx.PostForm("allergies")
	dob := ctx.PostForm("date_of_birth")
	gender := ctx.PostForm("gender")

	imagePath := c.handleImageUpload(ctx, userID)

	user, err := c.service.UpdateProfile(uint(userID.(float64)), fullName, phone, imagePath, bloodType, height, weight, allergies, dob, gender)
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", user))
}

func (c *ProfileController) handleImageUpload(ctx *gin.Context, userID any) string {
	file, err := ctx.FormFile("image")
	if err != nil {
		return ""
	}
	filename := fmt.Sprintf("%d_%d_%s", int(userID.(float64)), time.Now().Unix(), file.Filename)
	dst := "uploads/" + filename
	if err := ctx.SaveUploadedFile(file, dst); err == nil {
		return "http://127.0.0.1:8080/uploads/" + filename
	}
	return ""
}

func (c *ProfileController) AddFamilyMember(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}

	var req models.FamilyMember
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	req.UserID = uint(userID.(float64))
	if err := c.service.AddFamilyMember(&req); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", gin.H{
		"message":       "Family member added successfully",
		"family_member": req,
	}))
}

func (c *ProfileController) DeleteFamilyMember(ctx *gin.Context) {
	_, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(utils.NewValidationError("Invalid ID"))
		return
	}

	if err := c.service.DeleteFamilyMember(uint(id)); err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", gin.H{
		"message": "Family member deleted successfully",
	}))
}
