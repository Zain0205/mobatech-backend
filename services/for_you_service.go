package services

import (
	"backend/constants"
	"backend/repositories"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Article struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	ReadTime string `json:"readTime"`
	Content  string `json:"content"`
}

type ForYouService interface {
	GenerateRecommendations(ctx context.Context, userID string) ([]Article, error)
}

type forYouService struct {
	chatRepo repositories.ChatRepository
}

func NewForYouService(chatRepo repositories.ChatRepository) ForYouService {
	return &forYouService{chatRepo}
}

func (s *forYouService) GenerateRecommendations(ctx context.Context, userID string) ([]Article, error) {
	chatContext := s.extractUserContext(userID)
	prompt := s.buildPrompt(chatContext)

	jsonStr, err := s.fetchGeminiResponse(ctx, prompt)
	if err != nil {
		return s.getFallbackArticles(), nil
	}

	var articles []Article
	if err := json.Unmarshal([]byte(jsonStr), &articles); err != nil {
		return s.getFallbackArticles(), nil
	}

	return articles, nil
}

func (s *forYouService) extractUserContext(userID string) string {
	sessions, err := s.chatRepo.GetUserSessions(userID)
	if err != nil || len(sessions) == 0 {
		return "Pengguna belum pernah konsultasi. Berikan saran kesehatan umum."
	}

	chatContext := ""
	for i, session := range sessions {
		if i >= 3 {
			break
		}
		msgs, err := s.chatRepo.GetSessionMessages(session.ID)
		if err == nil {
			for _, msg := range msgs {
				if msg.Role == "user" {
					chatContext += "- " + msg.Content + "\n"
				}
			}
		}
	}

	if chatContext == "" {
		return "Pengguna belum pernah konsultasi. Berikan saran kesehatan umum."
	}
	return chatContext
}

func (s *forYouService) buildPrompt(context string) string {
	return fmt.Sprintf(constants.GeminiForYouPrompt, context)
}

func (s *forYouService) fetchGeminiResponse(ctx context.Context, prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	model.ResponseMIMEType = "application/json"

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	var jsonStr string
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				jsonStr += string(text)
			}
		}
	}
	return jsonStr, nil
}

func (s *forYouService) getFallbackArticles() []Article {
	return []Article{
		{Title: "Menjaga Pola Makan", Category: "Kesehatan", ReadTime: "12 Menit", Content: "Menjaga pola makan sehat sangat penting bagi keseimbangan tubuh..."},
		{Title: "Olahraga Ringan", Category: "Kebugaran", ReadTime: "8 Menit", Content: "Peregangan ringan selama 15 menit setiap pagi dapat melancarkan peredaran darah..."},
		{Title: "Tidur yang Cukup", Category: "Gaya Hidup", ReadTime: "5 Menit", Content: "Tidur 7-8 jam sangat dianjurkan untuk proses pemulihan..."},
		{Title: "Mengatasi Stres", Category: "Mental", ReadTime: "10 Menit", Content: "Manajemen stres dengan teknik pernapasan 4-7-8..."},
	}
}
