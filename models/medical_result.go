package models

import (
	"time"

	"gorm.io/gorm"
)

type MedicalResult struct {
	gorm.Model
	UserID        uint      `json:"user_id"`
	AppointmentID uint      `json:"appointment_id"`
	DoctorName    string    `json:"doctor_name"`
	TestType      string    `json:"test_type"`
	TestName      string    `json:"test_name"`
	Result        string    `json:"result" gorm:"type:text"`
	Notes         string    `json:"notes" gorm:"type:text"`
	FileURL       string    `json:"file_url"`
	ResultDate    time.Time `json:"result_date"`
}
