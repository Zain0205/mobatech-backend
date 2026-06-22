package models

import (
	"time"

	"gorm.io/gorm"
)

type Reminder struct {
	gorm.Model
	UserID        uint      `json:"user_id"`
	AppointmentID uint      `json:"appointment_id"`
	Title         string    `json:"title"`
	Message       string    `json:"message" gorm:"type:text"`
	ReminderDate  time.Time `json:"reminder_date"`
	IsRead        bool      `json:"is_read" gorm:"default:false"`
	Type          string    `json:"type"`
}
