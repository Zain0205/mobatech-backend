package models

import "gorm.io/gorm"

type Prescription struct {
	gorm.Model
	UserID           uint               `json:"user_id"`
	DoctorName       string             `json:"doctor_name"`
	DoctorSpecialty  string             `json:"doctor_specialty"`
	Diagnosis        string             `json:"diagnosis" gorm:"type:text"`
	PrescriptionDate string             `json:"prescription_date"`
	ExpiryDate       string             `json:"expiry_date"`
	Status           string             `json:"status"` // Active, Redeemed, Expired
	Notes            string             `json:"notes" gorm:"type:text"`
	Items            []PrescriptionItem `json:"items" gorm:"foreignKey:PrescriptionID"`
}

type PrescriptionItem struct {
	gorm.Model
	PrescriptionID    uint     `json:"prescription_id"`
	MedicineID        uint     `json:"medicine_id"`
	Medicine          Medicine `json:"medicine" gorm:"foreignKey:MedicineID"`
	DosageInstruction string   `json:"dosage_instruction"` // e.g., "3x1 sehari sesudah makan"
	Duration          string   `json:"duration"`           // e.g., "7 hari"
	Quantity          int      `json:"quantity"`
	Notes             string   `json:"notes"`
}
