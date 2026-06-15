package repositories

import (
	"backend/models"
	"gorm.io/gorm"
)

type ChatRepository interface {
	CreateSession(session *models.ChatSession) error
	GetUserSessions(userID string) ([]models.ChatSession, error)
	GetSession(sessionID uint) (*models.ChatSession, error)
	AddMessage(message *models.ChatMessage) error
	GetSessionMessages(sessionID uint) ([]models.ChatMessage, error)
	DeleteSession(userID uint, sessionID uint) error
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db}
}

func (r *chatRepository) CreateSession(session *models.ChatSession) error {
	return r.db.Create(session).Error
}

func (r *chatRepository) GetUserSessions(userID string) ([]models.ChatSession, error) {
	var sessions []models.ChatSession
	err := r.db.Where("user_id = ?", userID).Order("updated_at desc").Find(&sessions).Error
	return sessions, err
}

func (r *chatRepository) GetSession(sessionID uint) (*models.ChatSession, error) {
	var session models.ChatSession
	err := r.db.Preload("Messages").First(&session, sessionID).Error
	return &session, err
}

func (r *chatRepository) AddMessage(message *models.ChatMessage) error {
	return r.db.Create(message).Error
}

func (r *chatRepository) GetSessionMessages(sessionID uint) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Where("session_id = ?", sessionID).Order("created_at asc").Find(&messages).Error
	return messages, err
}

func (r *chatRepository) DeleteSession(userID uint, sessionID uint) error {
	r.db.Where("chat_session_id = ?", sessionID).Delete(&models.ChatMessage{})
	return r.db.Where("id = ? AND user_id = ?", sessionID, userID).Delete(&models.ChatSession{}).Error
}
