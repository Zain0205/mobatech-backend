package controllers

import (
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	user, err := c.service.Register(req.FullName, req.Email, req.PhoneNumber, req.Password)
	if err != nil {
		ctx.Error(utils.NewAppError(utils.ErrConflict, http.StatusConflict, "Email sudah terdaftar atau terjadi kesalahan server."))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", user))
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	token, user, err := c.service.Login(req.Email, req.Password)
	if err != nil {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Email atau kata sandi salah."))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", gin.H{"token": token, "user": user}))
}

func (c *AuthController) Me(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}

	user, err := c.service.GetUser(uint(userID.(float64)))
	if err != nil {
		ctx.Error(utils.NewAppError(utils.ErrNotFound, http.StatusNotFound, "User not found"))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", user))
}
