package services

import (
	"backend/models"
	"backend/repositories"
)

type EmergencyService interface {
	CreateRequest(req *models.EmergencyRequest) error
	GetByID(id uint) (*models.EmergencyRequest, error)
	GetHistoryByUser(userID uint) ([]models.EmergencyRequest, error)
	GetAllRequests() ([]models.EmergencyRequest, error)
	UpdateStatus(id uint, status string) error
	UpdateTracking(id uint, ambulanceLat, ambulanceLng float64, estimatedMinutes int, status string) error
}

type emergencyService struct {
	repo repositories.EmergencyRepository
}

func NewEmergencyService(repo repositories.EmergencyRepository) EmergencyService {
	return &emergencyService{repo}
}

func (s *emergencyService) CreateRequest(req *models.EmergencyRequest) error {
	req.Status = "Pending"
	return s.repo.Create(req)
}

func (s *emergencyService) GetByID(id uint) (*models.EmergencyRequest, error) {
	return s.repo.GetByID(id)
}

func (s *emergencyService) GetHistoryByUser(userID uint) ([]models.EmergencyRequest, error) {
	return s.repo.GetByUserID(userID)
}

func (s *emergencyService) GetAllRequests() ([]models.EmergencyRequest, error) {
	return s.repo.GetAll()
}

func (s *emergencyService) UpdateStatus(id uint, status string) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *emergencyService) UpdateTracking(id uint, ambulanceLat, ambulanceLng float64, estimatedMinutes int, status string) error {
	return s.repo.UpdateTracking(id, ambulanceLat, ambulanceLng, estimatedMinutes, status)
}
