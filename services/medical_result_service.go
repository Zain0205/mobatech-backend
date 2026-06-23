package services

import (
	"backend/models"
	"backend/repositories"
	"errors"
)

type MedicalResultService interface {
	GetAllMedicalResults() ([]models.MedicalResult, error)
	GetUserMedicalResults(userID uint) ([]models.MedicalResult, error)
	GetMedicalResultByID(id uint) (*models.MedicalResult, error)
	CreateMedicalResult(req *models.MedicalResult) (*models.MedicalResult, error)
	UpdateMedicalResult(req *models.MedicalResult) (*models.MedicalResult, error)
	DeleteMedicalResult(id uint) error
}

type medicalResultService struct {
	repo repositories.MedicalResultRepository
}

func NewMedicalResultService(repo repositories.MedicalResultRepository) MedicalResultService {
	return &medicalResultService{repo}
}

func (s *medicalResultService) GetAllMedicalResults() ([]models.MedicalResult, error) {
	return s.repo.FindAll()
}

func (s *medicalResultService) GetUserMedicalResults(userID uint) ([]models.MedicalResult, error) {
	return s.repo.FindByUserID(userID)
}

func (s *medicalResultService) GetMedicalResultByID(id uint) (*models.MedicalResult, error) {
	return s.repo.FindByID(id)
}

func (s *medicalResultService) CreateMedicalResult(req *models.MedicalResult) (*models.MedicalResult, error) {
	err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (s *medicalResultService) UpdateMedicalResult(req *models.MedicalResult) (*models.MedicalResult, error) {
	_, err := s.repo.FindByID(req.ID)
	if err != nil {
		return nil, errors.New("medical result not found")
	}

	err = s.repo.Update(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (s *medicalResultService) DeleteMedicalResult(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("medical result not found")
	}
	return s.repo.Delete(id)
}
