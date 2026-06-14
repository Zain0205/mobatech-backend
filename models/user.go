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
}
