package services

import (
	"backend/models"
	"backend/repositories"
	"errors"
)

type ReminderService interface {
	GetAllReminders() ([]models.Reminder, error)
	GetUserReminders(userID uint) ([]models.Reminder, error)
	GetUnreadCount(userID uint) (int64, error)
	CreateReminder(req *models.Reminder) (*models.Reminder, error)
	MarkAsRead(id uint, userID uint) error
	DeleteReminder(id uint) error
}

type reminderService struct {
	repo repositories.ReminderRepository
}

func NewReminderService(repo repositories.ReminderRepository) ReminderService {
	return &reminderService{repo}
}

func (s *reminderService) GetAllReminders() ([]models.Reminder, error) {
	return s.repo.FindAll()
}

func (s *reminderService) GetUserReminders(userID uint) ([]models.Reminder, error) {
	return s.repo.FindByUserID(userID)
}

func (s *reminderService) GetUnreadCount(userID uint) (int64, error) {
	return s.repo.FindUnreadCountByUserID(userID)
}

func (s *reminderService) CreateReminder(req *models.Reminder) (*models.Reminder, error) {
	err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (s *reminderService) MarkAsRead(id uint, userID uint) error {
	reminder, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("reminder not found")
	}
	
	if reminder.UserID != userID {
		return errors.New("unauthorized to update this reminder")
	}
	
	reminder.IsRead = true
	return s.repo.Update(reminder)
}

func (s *reminderService) DeleteReminder(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("reminder not found")
	}
	return s.repo.Delete(id)
}
