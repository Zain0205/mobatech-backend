package models

import "gorm.io/gorm"

type FamilyMember struct {
	gorm.Model
	UserID       uint   `json:"user_id"`
	FullName     string `json:"full_name"`
	Relationship string `json:"relationship"`
	DateOfBirth  string `json:"date_of_birth"`
	Gender       string `json:"gender"`
}
