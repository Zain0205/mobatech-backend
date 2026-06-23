package services

import (
	"backend/models"
	"backend/repositories"
	"errors"
)

type AppointmentService interface {
	GetAllAppointments() ([]models.Appointment, error)
	GetUserAppointments(userID uint) ([]models.Appointment, error)
	BookAppointment(userID uint, req *models.Appointment) (*models.Appointment, error)
	CancelAppointment(id uint, userID uint, isAdmin bool) error
	ApproveAppointment(id uint) error
}

type appointmentService struct {
	appointmentRepo repositories.AppointmentRepository
	scheduleRepo    repositories.ScheduleRepository
}

func NewAppointmentService(appointmentRepo repositories.AppointmentRepository, scheduleRepo repositories.ScheduleRepository) AppointmentService {
	return &appointmentService{appointmentRepo, scheduleRepo}
}

func (s *appointmentService) GetAllAppointments() ([]models.Appointment, error) {
	return s.appointmentRepo.FindAll()
}

func (s *appointmentService) GetUserAppointments(userID uint) ([]models.Appointment, error) {
	return s.appointmentRepo.FindByUserID(userID)
}

func (s *appointmentService) BookAppointment(userID uint, req *models.Appointment) (*models.Appointment, error) {
	// Verify schedule exists and has quota
	schedule, err := s.scheduleRepo.FindByID(req.DoctorScheduleID)
	if err != nil {
		return nil, errors.New("schedule not found")
	}

	if !schedule.IsAvailable || schedule.Booked >= schedule.Quota {
		return nil, errors.New("schedule is full or not available")
	}

	// Update schedule booked count
	schedule.Booked += 1
	err = s.scheduleRepo.Update(schedule)
	if err != nil {
		return nil, err
	}

	req.DoctorID = schedule.DoctorID
	req.UserID = userID
	req.Status = "pending"

	err = s.appointmentRepo.Create(req)
	if err != nil {
		s.rollbackScheduleBooking(schedule)
		return nil, err
	}

	return s.appointmentRepo.FindByID(req.ID)
}

func (s *appointmentService) CancelAppointment(id uint, userID uint, isAdmin bool) error {
	appointment, err := s.appointmentRepo.FindByID(id)
	if err != nil {
		return errors.New("appointment not found")
	}

	if !isAdmin && appointment.UserID != userID {
		return errors.New("unauthorized to cancel this appointment")
	}

	if appointment.Status == "cancelled" || appointment.Status == "completed" {
		return errors.New("cannot cancel an already cancelled or completed appointment")
	}

	appointment.Status = "cancelled"
	err = s.appointmentRepo.Update(appointment)
	if err != nil {
		return err
	}

	// Release schedule slot
	schedule, err := s.scheduleRepo.FindByID(appointment.DoctorScheduleID)
	if err == nil && schedule.Booked > 0 {
		schedule.Booked -= 1
		s.scheduleRepo.Update(schedule)
	}

	return nil
}

func (s *appointmentService) ApproveAppointment(id uint) error {
	appointment, err := s.appointmentRepo.FindByID(id)
	if err != nil {
		return errors.New("appointment not found")
	}

	if appointment.Status != "pending" {
		return errors.New("can only approve pending appointments")
	}

	appointment.Status = "approved"
	return s.appointmentRepo.Update(appointment)
}

func (s *appointmentService) rollbackScheduleBooking(schedule *models.DoctorSchedule) {
	schedule.Booked -= 1
	s.scheduleRepo.Update(schedule)
}
