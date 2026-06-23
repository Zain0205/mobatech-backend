package controllers

import (
	"backend/services"
	"backend/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ForYouController struct {
	service services.ForYouService
}

func NewForYouController(service services.ForYouService) *ForYouController {
	return &ForYouController{service}
}

func (ctrl *ForYouController) GetRecommendations(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "UNAUTHENTICATED"))
		return
	}

	userID := fmt.Sprintf("%v", userIDStr)
	articles, err := ctrl.service.GenerateRecommendations(c.Request.Context(), userID)
	if err != nil {
		c.Error(utils.NewInternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", articles))
}
