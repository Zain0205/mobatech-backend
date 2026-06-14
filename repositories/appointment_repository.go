package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type AppointmentRepository interface {
	FindAll() ([]models.Appointment, error)
	FindByUserID(userID uint) ([]models.Appointment, error)
	FindByID(id uint) (*models.Appointment, error)
	Create(appointment *models.Appointment) error
	Update(appointment *models.Appointment) error
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db}
}

func (r *appointmentRepository) FindAll() ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.Preload("User").Preload("Doctor").Preload("Schedule").
		Order("created_at desc").
		Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) FindByUserID(userID uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.Preload("Doctor").Preload("Schedule").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) FindByID(id uint) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.db.Preload("User").Preload("Doctor").Preload("Schedule").First(&appointment, id).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *appointmentRepository) Create(appointment *models.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *appointmentRepository) Update(appointment *models.Appointment) error {
	return r.db.Save(appointment).Error
}
