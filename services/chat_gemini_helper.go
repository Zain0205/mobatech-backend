package services

import (
	"backend/constants"
	"backend/models"
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func (s *chatService) setupGemini(ctx context.Context) (*genai.GenerativeModel, *genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create gemini client: %v", err)
	}

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(constants.GeminiSystemPrompt)},
	}
	return model, client, nil
}

func (s *chatService) populateHistory(cs *genai.ChatSession, historyMsg []models.ChatMessage, lastUserMsg string) {
	for _, msg := range historyMsg {
		if msg.Role == "user" && msg.Content == lastUserMsg {
			continue
		}
		role := msg.Role
		if role != "model" && role != "user" {
			role = "user"
		}
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{genai.Text(msg.Content)},
			Role:  role,
		})
	}
}

func (s *chatService) processStream(iter *genai.GenerateContentResponseIterator, outChan chan<- string, errChan chan<- error, sessionID uint) {
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

	s.repo.AddMessage(&models.ChatMessage{
		SessionID: sessionID,
		Role:      "model",
		Content:   fullResponse,
	})
}
