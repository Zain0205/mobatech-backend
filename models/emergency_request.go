package models

import "gorm.io/gorm"

type EmergencyRequest struct {
	gorm.Model
	UserID           uint    `json:"user_id"`
	PatientName      string  `json:"patient_name"`
	Condition        string  `json:"condition"`
	Location         string  `json:"location"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	PhoneNumber      string  `json:"phone_number"`
	Status           string  `json:"status"` // e.g. "Pending", "Dispatched", "Arrived", "Resolved"
	AmbulanceLat     float64 `json:"ambulance_lat"`
	AmbulanceLng     float64 `json:"ambulance_lng"`
	EstimatedMinutes int     `json:"estimated_minutes"`
}
