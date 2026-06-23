package controllers

import (
	"backend/services"
	"backend/utils"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	service services.ChatService
}

func NewChatController(service services.ChatService) *ChatController {
	return &ChatController{service}
}

func (c *ChatController) CreateSession(ctx *gin.Context) {
	var req struct {
		Title string `json:"title"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	userIDStr, _ := ctx.Get("user_id")
	userID := fmt.Sprintf("%v", userIDStr)

	session, err := c.service.CreateSession(userID, req.Title)
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", session))
}

func (c *ChatController) GetUserSessions(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}

	userID := fmt.Sprintf("%v", userIDStr)
	sessions, err := c.service.GetUserSessions(userID)
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", sessions))
}

func (c *ChatController) GetSessionMessages(ctx *gin.Context) {
	sessionIDStr := ctx.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		ctx.Error(utils.NewValidationError("invalid session id"))
		return
	}

	messages, err := c.service.GetSessionMessages(uint(sessionID))
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", messages))
}

func (c *ChatController) DeleteSession(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Unauthorized"))
		return
	}

	sessionIDStr := ctx.Param("id")
	sessionID, _ := strconv.Atoi(sessionIDStr)

	err := c.service.DeleteSession(uint(userID.(float64)), uint(sessionID))
	if err != nil {
		ctx.Error(utils.NewInternalError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildSuccess("OK", "Success", nil))
}

func (c *ChatController) StreamChat(ctx *gin.Context) {
	sessionIDStr := ctx.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		ctx.Error(utils.NewValidationError("invalid session id"))
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Error(utils.NewValidationError(err.Error()))
		return
	}

	outChan := make(chan string)
	errChan := make(chan error)

	// Context from Request is passed to manage cancellation
	go c.service.StreamChat(ctx.Request.Context(), uint(sessionID), req.Message, outChan, errChan)

	c.handleStream(ctx, outChan, errChan)
}

func (c *ChatController) handleStream(ctx *gin.Context, outChan <-chan string, errChan <-chan error) {
	ctx.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-outChan:
			if !ok {
				return false
			}
			ctx.SSEvent("message", gin.H{"text": msg})
			return true
		case err, ok := <-errChan:
			if !ok {
				return false
			}
			ctx.SSEvent("error", err.Error())
			return false
		case <-ctx.Request.Context().Done():
			return false
		}
	})
}
