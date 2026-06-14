package models

import "gorm.io/gorm"

type Appointment struct {
	gorm.Model
	UserID           uint            `json:"user_id"`
	User             *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DoctorID         uint            `json:"doctor_id"`
	Doctor           *Doctor         `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	DoctorScheduleID uint            `json:"doctor_schedule_id"`
	Schedule         *DoctorSchedule `json:"schedule,omitempty" gorm:"foreignKey:DoctorScheduleID"`
	Status           string          `json:"status" gorm:"default:'pending'"` // pending, approved, cancelled, completed
	Notes            string          `json:"notes"`
}
