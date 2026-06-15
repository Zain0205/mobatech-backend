package controllers

import (
	"backend/services"
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
		Title  string `json:"title"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := ctx.Get("user_id")
	userID := fmt.Sprintf("%v", userIDStr)

	session, err := c.service.CreateSession(userID, req.Title)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, session)
}

func (c *ChatController) GetUserSessions(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := fmt.Sprintf("%v", userIDStr)
	sessions, err := c.service.GetUserSessions(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sessions)
}

func (c *ChatController) GetSessionMessages(ctx *gin.Context) {
	sessionIDStr := ctx.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	messages, err := c.service.GetSessionMessages(uint(sessionID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

func (c *ChatController) DeleteSession(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sessionIDStr := ctx.Param("id")
	sessionID, _ := strconv.Atoi(sessionIDStr)

	err := c.service.DeleteSession(uint(userID.(float64)), uint(sessionID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Session deleted"})
}

func (c *ChatController) StreamChat(ctx *gin.Context) {
	sessionIDStr := ctx.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outChan := make(chan string)
	errChan := make(chan error)

	// Context from Request is passed to manage cancellation
	go c.service.StreamChat(ctx.Request.Context(), uint(sessionID), req.Message, outChan, errChan)

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
