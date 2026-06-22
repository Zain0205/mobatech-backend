package controllers

import (
	"backend/models"
	"backend/services"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func formatValidationError(err error) gin.H {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string][]string)
		for _, e := range validationErrs {
			field := e.Field()
			var msg string
			switch e.Tag() {
			case "required":
				msg = "Kolom ini wajib diisi"
			case "email":
				msg = "Format email tidak valid"
			case "min":
				msg = "Minimal " + e.Param() + " karakter"
			default:
				msg = "Input tidak valid"
			}
			errors[field] = append(errors[field], msg)
		}
		return gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Validasi gagal. Silakan periksa kembali input Anda.",
			"errors":  errors,
		}
	}
	return gin.H{
		"success": false,
		"code":    "BAD_REQUEST",
		"message": "Format permintaan tidak valid",
		"error":   err.Error(),
	}
}

type AuthController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{service}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req struct {
		FullName    string `json:"full_name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, formatValidationError(err))
		return
	}

	user, err := c.service.Register(req.FullName, req.Email, req.PhoneNumber, req.Password)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"code":    "REGISTER_ERROR",
			"message": "Gagal mendaftar. " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, formatValidationError(err))
		return
	}

	token, user, err := c.service.Login(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "UNAUTHENTICATED",
			"message": "Email atau kata sandi salah.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

func (c *AuthController) Me(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := c.service.GetUser(uint(userID.(float64)))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fullName := ctx.PostForm("full_name")
	phone := ctx.PostForm("phone_number")

	bloodType := ctx.PostForm("blood_type")
	heightStr := ctx.PostForm("height")
	weightStr := ctx.PostForm("weight")
	allergies := ctx.PostForm("allergies")
	dob := ctx.PostForm("date_of_birth")
	gender := ctx.PostForm("gender")

	height, _ := strconv.Atoi(heightStr)
	weight, _ := strconv.Atoi(weightStr)

	file, err := ctx.FormFile("image")
	imagePath := ""
	if err == nil {
		filename := fmt.Sprintf("%d_%d_%s", int(userID.(float64)), time.Now().Unix(), file.Filename)
		dst := "uploads/" + filename
		if err := ctx.SaveUploadedFile(file, dst); err == nil {
			imagePath = "http://127.0.0.1:8080/uploads/" + filename
		}
	}

	user, err := c.service.UpdateProfile(uint(userID.(float64)), fullName, phone, imagePath, bloodType, height, weight, allergies, dob, gender)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

func (c *AuthController) AddFamilyMember(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.FamilyMember
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = uint(userID.(float64))
	if err := c.service.AddFamilyMember(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Family member added successfully",
		"family_member": req,
	})
}

func (c *AuthController) DeleteFamilyMember(ctx *gin.Context) {
	_, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.DeleteFamilyMember(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Family member deleted successfully",
	})
}
