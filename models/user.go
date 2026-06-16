package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName    string `json:"full_name"`
	Email       string `json:"email" gorm:"unique"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"-"`
	ImageURL    string `json:"image_url"`
	BloodType   string `json:"blood_type"`
	Height      int    `json:"height"`
	Weight      int    `json:"weight"`
	Allergies   string `json:"allergies"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	FamilyMembers []FamilyMember `json:"family_members" gorm:"foreignKey:UserID"`
}
