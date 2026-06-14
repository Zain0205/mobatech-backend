package models

import "gorm.io/gorm"

type Medicine struct {
	gorm.Model
	CategoryID           uint             `json:"category_id"`
	Category             MedicineCategory `json:"category" gorm:"foreignKey:CategoryID"`
	Name                 string           `json:"name"`
	GenericName          string           `json:"generic_name"`
	Description          string           `json:"description" gorm:"type:text"`
	Dosage               string           `json:"dosage"`
	Unit                 string           `json:"unit"`
	Price                float64          `json:"price"`
	Stock                int              `json:"stock"`
	RequiresPrescription bool             `json:"requires_prescription"`
	Manufacturer         string           `json:"manufacturer"`
	ImageURL             string           `json:"image_url"`
}
