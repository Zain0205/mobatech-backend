package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type MedicalResultRepository interface {
	FindAll() ([]models.MedicalResult, error)
	FindByUserID(userID uint) ([]models.MedicalResult, error)
	FindByID(id uint) (*models.MedicalResult, error)
	Create(medicalResult *models.MedicalResult) error
	Update(medicalResult *models.MedicalResult) error
	Delete(id uint) error
}

type medicalResultRepository struct {
	db *gorm.DB
}

func NewMedicalResultRepository(db *gorm.DB) MedicalResultRepository {
	return &medicalResultRepository{db}
}

func (r *medicalResultRepository) FindAll() ([]models.MedicalResult, error) {
	var results []models.MedicalResult
	err := r.db.Order("created_at desc").Find(&results).Error
	return results, err
}

func (r *medicalResultRepository) FindByUserID(userID uint) ([]models.MedicalResult, error) {
	var results []models.MedicalResult
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&results).Error
	return results, err
}

func (r *medicalResultRepository) FindByID(id uint) (*models.MedicalResult, error) {
	var result models.MedicalResult
	err := r.db.First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *medicalResultRepository) Create(medicalResult *models.MedicalResult) error {
	return r.db.Create(medicalResult).Error
}

func (r *medicalResultRepository) Update(medicalResult *models.MedicalResult) error {
	return r.db.Save(medicalResult).Error
}

func (r *medicalResultRepository) Delete(id uint) error {
	return r.db.Delete(&models.MedicalResult{}, id).Error
}
