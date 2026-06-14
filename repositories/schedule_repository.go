package repositories

import (
	"backend/models"
	"time"

	"gorm.io/gorm"
)

type ScheduleRepository interface {
	FindByDoctorID(doctorID uint, fromDate time.Time) ([]models.DoctorSchedule, error)
	FindByID(id uint) (*models.DoctorSchedule, error)
	Create(schedule *models.DoctorSchedule) error
	Update(schedule *models.DoctorSchedule) error
	Delete(id uint) error
}

type scheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{db}
}

func (r *scheduleRepository) FindByDoctorID(doctorID uint, fromDate time.Time) ([]models.DoctorSchedule, error) {
	var schedules []models.DoctorSchedule
	err := r.db.Preload("Doctor").
		Where("doctor_id = ? AND date >= ? AND is_available = ?", doctorID, fromDate, true).
		Order("date asc, start_time asc").
		Find(&schedules).Error
	return schedules, err
}

func (r *scheduleRepository) FindByID(id uint) (*models.DoctorSchedule, error) {
	var schedule models.DoctorSchedule
	err := r.db.Preload("Doctor").First(&schedule, id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) Create(schedule *models.DoctorSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *scheduleRepository) Update(schedule *models.DoctorSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *scheduleRepository) Delete(id uint) error {
	return r.db.Model(&models.DoctorSchedule{}).Where("id = ?", id).Update("is_available", false).Error
}
