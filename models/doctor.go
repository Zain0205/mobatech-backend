package models

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	UserID         *uint  `json:"user_id,omitempty"` // Optional user account link
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
	ContactInfo    string `json:"contact_info"`
	Description    string `json:"description"`
	ImageURL       string `json:"image_url"`
	IsActive       bool   `json:"is_active" gorm:"default:true"`
}
