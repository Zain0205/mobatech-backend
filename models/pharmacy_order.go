package models

import "gorm.io/gorm"

type PharmacyOrder struct {
	gorm.Model
	UserID          uint                `json:"user_id"`
	PrescriptionID  *uint               `json:"prescription_id"` // nullable for OTC
	OrderNumber     string              `json:"order_number" gorm:"uniqueIndex"`
	Status          string              `json:"status"` // Pending, Verifying, Processing, Ready, Completed, Cancelled
	TotalPrice      float64             `json:"total_price"`
	PaymentMethod   string              `json:"payment_method"` // Cash, BPJS, Insurance, Transfer
	PaymentStatus   string              `json:"payment_status"` // Unpaid, Paid, Refunded
	PickupMethod    string              `json:"pickup_method"`  // Counter, Delivery
	DeliveryAddress string              `json:"delivery_address" gorm:"type:text"`
	Notes           string              `json:"notes" gorm:"type:text"`
	Items           []PharmacyOrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type PharmacyOrderItem struct {
	gorm.Model
	OrderID    uint     `json:"order_id"`
	MedicineID uint     `json:"medicine_id"`
	Medicine   Medicine `json:"medicine" gorm:"foreignKey:MedicineID"`
	Quantity   int      `json:"quantity"`
	Price      float64  `json:"price"`
	Subtotal   float64  `json:"subtotal"`
	Notes      string   `json:"notes"`
}
