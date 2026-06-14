package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type DoctorRepository interface {
	FindAll(specialization string) ([]models.Doctor, error)
	FindByID(id uint) (*models.Doctor, error)
	Create(doctor *models.Doctor) error
	Update(doctor *models.Doctor) error
	Delete(id uint) error
}

type doctorRepository struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepository{db}
}

func (r *doctorRepository) FindAll(specialization string) ([]models.Doctor, error) {
	var doctors []models.Doctor
	query := r.db.Where("is_active = ?", true)
	if specialization != "" {
		query = query.Where("specialization = ?", specialization)
	}
	err := query.Find(&doctors).Error
	return doctors, err
}

func (r *doctorRepository) FindByID(id uint) (*models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.First(&doctor, id).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorRepository) Create(doctor *models.Doctor) error {
	return r.db.Create(doctor).Error
}

func (r *doctorRepository) Update(doctor *models.Doctor) error {
	return r.db.Save(doctor).Error
}

func (r *doctorRepository) Delete(id uint) error {
	return r.db.Model(&models.Doctor{}).Where("id = ?", id).Update("is_active", false).Error
}
