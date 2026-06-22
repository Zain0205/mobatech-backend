package services

import (
	"backend/models"
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
	sessions, err := s.chatRepo.GetUserSessions(userID)
	if err != nil {
		return nil, err
	}

	chatContext := ""
	for i, session := range sessions {
		if i >= 3 {
			break
		}
		messages, err := s.chatRepo.GetSessionMessages(session.ID)
		if err == nil {
			for _, msg := range messages {
				if msg.Role == "user" {
					chatContext += "- " + msg.Content + "\n"
				}
			}
		}
	}

	if chatContext == "" {
		chatContext = "Pengguna belum pernah konsultasi. Berikan saran kesehatan umum."
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat genai client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	model.ResponseMIMEType = "application/json"
	
	prompt := fmt.Sprintf(`Anda adalah AI Asisten Medis profesional. Berdasarkan keluhan/pertanyaan pengguna ini:
%s

Hasilkan 4 rekomendasi artikel/tips kesehatan yang PALING RELEVAN dengan kondisi tersebut.
Format output HARUS JSON array, persis seperti ini:
[
  {
    "title": "string (Judul menarik)",
    "category": "string (Kategori medis)",
    "readTime": "string (misal '5 Menit baca')",
    "content": "string (ISI LENGKAP ARTIKEL. Berikan edukasi medis, tips, dan anjuran yang mendalam, sekitar 2-3 paragraf lengkap dengan formatting newline \n jika perlu)"
  }
]`, chatContext)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gagal generate: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("kosong")
	}

	var jsonStr string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			jsonStr += string(text)
		}
	}

	var articles []Article
	err = json.Unmarshal([]byte(jsonStr), &articles)
	if err != nil {
		return []Article{
			{Title: "Pentingnya Menjaga Pola Makan", Category: "Kesehatan", ReadTime: "12 Menit baca", Content: "Menjaga pola makan sehat sangat penting bagi keseimbangan tubuh dan pencegahan penyakit kronis. Nutrisi yang seimbang membantu sistem kekebalan tubuh bekerja maksimal..."},
			{Title: "Tips Olahraga Ringan di Rumah", Category: "Kebugaran", ReadTime: "8 Menit baca", Content: "Olahraga tidak harus di gym. Peregangan ringan selama 15 menit setiap pagi dapat melancarkan peredaran darah dan meningkatkan fokus kerja Anda seharian..."},
			{Title: "Pentingnya Tidur yang Cukup", Category: "Gaya Hidup", ReadTime: "5 Menit baca", Content: "Tidur 7-8 jam sangat dianjurkan untuk proses pemulihan sel-sel otak dan tubuh. Kurang tidur terbukti meningkatkan risiko obesitas dan penyakit kardiovaskular..."},
			{Title: "Cara Mengatasi Stres Kerja", Category: "Mental", ReadTime: "10 Menit baca", Content: "Manajemen stres dapat dilakukan dengan teknik pernapasan 4-7-8, meluangkan waktu hobi, dan membatasi screen time setelah jam kerja selesai..."},
		}, nil
	}

	return articles, nil
}
