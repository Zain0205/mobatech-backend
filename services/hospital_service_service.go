package services

import (
	"backend/models"
	"backend/repositories"
)

type HospitalServiceService interface {
	GetAll() ([]models.HospitalService, error)
	GetByID(id uint) (*models.HospitalService, error)
	Create(service *models.HospitalService) error
	Update(service *models.HospitalService) error
	Delete(id uint) error
}

type hospitalServiceService struct {
	repo repositories.HospitalServiceRepository
}

func NewHospitalServiceService(repo repositories.HospitalServiceRepository) HospitalServiceService {
	return &hospitalServiceService{repo}
}

func (s *hospitalServiceService) GetAll() ([]models.HospitalService, error) {
	return s.repo.GetAll()
}

func (s *hospitalServiceService) GetByID(id uint) (*models.HospitalService, error) {
	return s.repo.GetByID(id)
}

func (s *hospitalServiceService) Create(service *models.HospitalService) error {
	return s.repo.Create(service)
}

func (s *hospitalServiceService) Update(service *models.HospitalService) error {
	return s.repo.Update(service)
}

func (s *hospitalServiceService) Delete(id uint) error {
	return s.repo.Delete(id)
}
