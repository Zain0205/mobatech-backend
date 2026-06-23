package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type HospitalServiceRepository interface {
	GetAll() ([]models.HospitalService, error)
	GetByID(id uint) (*models.HospitalService, error)
	Create(service *models.HospitalService) error
	Update(service *models.HospitalService) error
	Delete(id uint) error
}

type hospitalServiceRepository struct {
	db *gorm.DB
}

func NewHospitalServiceRepository(db *gorm.DB) HospitalServiceRepository {
	return &hospitalServiceRepository{db}
}

func (r *hospitalServiceRepository) GetAll() ([]models.HospitalService, error) {
	var services []models.HospitalService
	err := r.db.Find(&services).Error
	return services, err
}

func (r *hospitalServiceRepository) GetByID(id uint) (*models.HospitalService, error) {
	var service models.HospitalService
	err := r.db.First(&service, id).Error
	return &service, err
}

func (r *hospitalServiceRepository) Create(service *models.HospitalService) error {
	return r.db.Create(service).Error
}

func (r *hospitalServiceRepository) Update(service *models.HospitalService) error {
	return r.db.Save(service).Error
}

func (r *hospitalServiceRepository) Delete(id uint) error {
	return r.db.Delete(&models.HospitalService{}, id).Error
}
