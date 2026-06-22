package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type ReminderRepository interface {
	FindAll() ([]models.Reminder, error)
	FindByUserID(userID uint) ([]models.Reminder, error)
	FindUnreadCountByUserID(userID uint) (int64, error)
	FindByID(id uint) (*models.Reminder, error)
	Create(reminder *models.Reminder) error
	Update(reminder *models.Reminder) error
	Delete(id uint) error
}

type reminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepository{db}
}

func (r *reminderRepository) FindAll() ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Order("created_at desc").Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) FindByUserID(userID uint) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Where("user_id = ?", userID).Order("reminder_date asc").Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) FindUnreadCountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Reminder{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (r *reminderRepository) FindByID(id uint) (*models.Reminder, error) {
	var reminder models.Reminder
	err := r.db.First(&reminder, id).Error
	if err != nil {
		return nil, err
	}
	return &reminder, nil
}

func (r *reminderRepository) Create(reminder *models.Reminder) error {
	return r.db.Create(reminder).Error
}

func (r *reminderRepository) Update(reminder *models.Reminder) error {
	return r.db.Save(reminder).Error
}

func (r *reminderRepository) Delete(id uint) error {
	return r.db.Delete(&models.Reminder{}, id).Error
}
