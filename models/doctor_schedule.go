package models

import (
	"time"

	"gorm.io/gorm"
)

type DoctorSchedule struct {
	gorm.Model
	DoctorID    uint      `json:"doctor_id"`
	Doctor      *Doctor   `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	Date        time.Time `json:"date" gorm:"type:date"`
	StartTime   string    `json:"start_time"` // e.g. "10:00"
	EndTime     string    `json:"end_time"`   // e.g. "12:00"
	Quota       int       `json:"quota"`
	Booked      int       `json:"booked" gorm:"default:0"`
	IsAvailable bool      `json:"is_available" gorm:"default:true"`
}
