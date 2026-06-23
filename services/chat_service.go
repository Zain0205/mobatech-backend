package services

import (
	"backend/models"
	"backend/repositories"
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

type ChatService interface {
	CreateSession(userID string, title string) (*models.ChatSession, error)
	GetUserSessions(userID string) ([]models.ChatSession, error)
	GetSessionMessages(sessionID uint) ([]models.ChatMessage, error)
	DeleteSession(userID uint, sessionID uint) error
	StreamChat(ctx context.Context, sessionID uint, userMessage string, outChan chan<- string, errChan chan<- error)
}

type chatService struct {
	repo repositories.ChatRepository
}

func NewChatService(repo repositories.ChatRepository) ChatService {
	return &chatService{repo}
}

func (s *chatService) CreateSession(userID string, title string) (*models.ChatSession, error) {
	session := &models.ChatSession{UserID: userID, Title: title}
	return session, s.repo.CreateSession(session)
}

func (s *chatService) GetUserSessions(userID string) ([]models.ChatSession, error) {
	return s.repo.GetUserSessions(userID)
}

func (s *chatService) GetSessionMessages(sessionID uint) ([]models.ChatMessage, error) {
	return s.repo.GetSessionMessages(sessionID)
}

func (s *chatService) DeleteSession(userID uint, sessionID uint) error {
	return s.repo.DeleteSession(userID, sessionID)
}

func (s *chatService) StreamChat(ctx context.Context, sessionID uint, userMessage string, outChan chan<- string, errChan chan<- error) {
	defer close(outChan)
	defer close(errChan)

	if err := s.saveUserMessage(sessionID, userMessage); err != nil {
		errChan <- err
		return
	}

	historyMsg, err := s.repo.GetSessionMessages(sessionID)
	if err != nil {
		errChan <- fmt.Errorf("failed to get history: %v", err)
		return
	}

	model, client, err := s.setupGemini(ctx)
	if err != nil {
		errChan <- err
		return
	}
	defer client.Close()

	cs := model.StartChat()
	s.populateHistory(cs, historyMsg, userMessage)

	// RAG Integration
	ragQuery := s.buildRAGPrompt(userMessage)

	iter := cs.SendMessageStream(ctx, genai.Text(ragQuery))
	s.processStream(iter, outChan, errChan, sessionID)
}

func (s *chatService) saveUserMessage(sessionID uint, userMessage string) error {
	err := s.repo.AddMessage(&models.ChatMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   userMessage,
	})
	if err != nil {
		return fmt.Errorf("failed to save user message: %v", err)
	}
	return nil
}

