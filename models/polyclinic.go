package models

import "gorm.io/gorm"

type Polyclinic struct {
	gorm.Model
	Name        string               `json:"name"`
	Description string               `json:"description"`
	ImageURL    string               `json:"image_url"`
	IsActive    bool                 `json:"is_active" gorm:"default:true"`
	Schedules   []PolyclinicSchedule `json:"schedules,omitempty" gorm:"foreignKey:PolyclinicID"`
}

type PolyclinicSchedule struct {
	gorm.Model
	PolyclinicID uint   `json:"polyclinic_id"`
	DayOfWeek    string `json:"day_of_week"` // e.g. "Senin", "Selasa"
	StartTime    string `json:"start_time"`  // e.g. "08:00"
	EndTime      string `json:"end_time"`    // e.g. "16:00"
	IsAvailable  bool   `json:"is_available" gorm:"default:true"`
}
