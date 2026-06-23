package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID uint       `json:"user_id" gorm:"uniqueIndex"`
	Items  []CartItem `json:"items" gorm:"foreignKey:CartID"`
}

type CartItem struct {
	gorm.Model
	CartID     uint     `json:"cart_id"`
	MedicineID uint     `json:"medicine_id"`
	Medicine   Medicine `json:"medicine" gorm:"foreignKey:MedicineID"`
	Quantity   int      `json:"quantity"`
}
