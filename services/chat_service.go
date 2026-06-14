package services

import (
	"backend/models"
	"backend/repositories"
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type ChatService interface {
	CreateSession(userID string, title string) (*models.ChatSession, error)
	GetUserSessions(userID string) ([]models.ChatSession, error)
	GetSessionMessages(sessionID uint) ([]models.ChatMessage, error)
	StreamChat(ctx context.Context, sessionID uint, userMessage string, outChan chan<- string, errChan chan<- error)
}

type chatService struct {
	repo repositories.ChatRepository
}

func NewChatService(repo repositories.ChatRepository) ChatService {
	return &chatService{repo}
}

func (s *chatService) CreateSession(userID string, title string) (*models.ChatSession, error) {
	session := &models.ChatSession{
		UserID: userID,
		Title:  title,
	}
	err := s.repo.CreateSession(session)
	return session, err
}

func (s *chatService) GetUserSessions(userID string) ([]models.ChatSession, error) {
	return s.repo.GetUserSessions(userID)
}

func (s *chatService) GetSessionMessages(sessionID uint) ([]models.ChatMessage, error) {
	return s.repo.GetSessionMessages(sessionID)
}

func (s *chatService) StreamChat(ctx context.Context, sessionID uint, userMessage string, outChan chan<- string, errChan chan<- error) {
	defer close(outChan)
	defer close(errChan)

	// Save user message
	err := s.repo.AddMessage(&models.ChatMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   userMessage,
	})
	if err != nil {
		errChan <- fmt.Errorf("failed to save user message: %v", err)
		return
	}

	historyMsg, err := s.repo.GetSessionMessages(sessionID)
	if err != nil {
		errChan <- fmt.Errorf("failed to get history: %v", err)
		return
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		errChan <- fmt.Errorf("failed to create gemini client: %v", err)
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	cs := model.StartChat()

	// Populate history
	for _, msg := range historyMsg {
		// skip the message we just saved
		if msg.Role == "user" && msg.Content == userMessage {
			continue
		}
		
		role := msg.Role
		if role == "model" {
			cs.History = append(cs.History, &genai.Content{
				Parts: []genai.Part{genai.Text(msg.Content)},
				Role:  "model",
			})
		} else {
			cs.History = append(cs.History, &genai.Content{
				Parts: []genai.Part{genai.Text(msg.Content)},
				Role:  "user",
			})
		}
	}

	iter := cs.SendMessageStream(ctx, genai.Text(userMessage))
	var fullResponse string

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			errChan <- fmt.Errorf("stream error: %v", err)
			return
		}
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					if text, ok := part.(genai.Text); ok {
						fullResponse += string(text)
						outChan <- string(text)
					}
				}
			}
		}
	}

	// Save model response
	s.repo.AddMessage(&models.ChatMessage{
		SessionID: sessionID,
		Role:      "model",
		Content:   fullResponse,
	})
}
