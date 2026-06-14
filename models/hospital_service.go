package models

import "gorm.io/gorm"

type HospitalService struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
