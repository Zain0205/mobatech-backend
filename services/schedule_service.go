package services

import (
	"backend/models"
	"backend/repositories"
	"time"
)

type ScheduleService interface {
	GetDoctorSchedules(doctorID uint) ([]models.DoctorSchedule, error)
	CreateSchedule(schedule *models.DoctorSchedule) error
	UpdateSchedule(id uint, input *models.DoctorSchedule) (*models.DoctorSchedule, error)
	DeleteSchedule(id uint) error
}

type scheduleService struct {
	scheduleRepo repositories.ScheduleRepository
}

func NewScheduleService(scheduleRepo repositories.ScheduleRepository) ScheduleService {
	return &scheduleService{scheduleRepo}
}

func (s *scheduleService) GetDoctorSchedules(doctorID uint) ([]models.DoctorSchedule, error) {
	// Only show schedules from today onwards
	now := time.Now().Truncate(24 * time.Hour)
	return s.scheduleRepo.FindByDoctorID(doctorID, now)
}

func (s *scheduleService) CreateSchedule(schedule *models.DoctorSchedule) error {
	schedule.IsAvailable = true
	schedule.Booked = 0
	return s.scheduleRepo.Create(schedule)
}

func (s *scheduleService) UpdateSchedule(id uint, input *models.DoctorSchedule) (*models.DoctorSchedule, error) {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !input.Date.IsZero() {
		schedule.Date = input.Date
	}
	if input.StartTime != "" {
		schedule.StartTime = input.StartTime
	}
	if input.EndTime != "" {
		schedule.EndTime = input.EndTime
	}
	if input.Quota > 0 {
		schedule.Quota = input.Quota
	}
	
	err = s.scheduleRepo.Update(schedule)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleService) DeleteSchedule(id uint) error {
	return s.scheduleRepo.Delete(id)
}
