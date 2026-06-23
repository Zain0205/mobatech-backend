package repositories

import (
	"backend/models"
)

func (r *pharmacyRepository) GetPrescriptionsByUserID(userID uint) ([]models.Prescription, error) {
	var prescriptions []models.Prescription
	err := r.db.Preload("Items").Preload("Items.Medicine").
		Where("user_id = ?", userID).Order("created_at desc").Find(&prescriptions).Error
	return prescriptions, err
}

func (r *pharmacyRepository) GetPrescriptionByID(id uint) (*models.Prescription, error) {
	var prescription models.Prescription
	err := r.db.Preload("Items").Preload("Items.Medicine").First(&prescription, id).Error
	return &prescription, err
}

func (r *pharmacyRepository) GetAllPrescriptions() ([]models.Prescription, error) {
	var prescriptions []models.Prescription
	err := r.db.Preload("Items").Preload("Items.Medicine").Order("created_at desc").Find(&prescriptions).Error
	return prescriptions, err
}

func (r *pharmacyRepository) CreatePrescription(p *models.Prescription) error {
	return r.db.Create(p).Error
}

func (r *pharmacyRepository) UpdatePrescriptionStatus(id uint, status string) error {
	return r.db.Model(&models.Prescription{}).Where("id = ?", id).Update("status", status).Error
}

func (r *pharmacyRepository) GetOrdersByUserID(userID uint) ([]models.PharmacyOrder, error) {
	var orders []models.PharmacyOrder
	err := r.db.Preload("Items").Preload("Items.Medicine").
		Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *pharmacyRepository) GetOrderByID(id uint) (*models.PharmacyOrder, error) {
	var order models.PharmacyOrder
	err := r.db.Preload("Items").Preload("Items.Medicine").First(&order, id).Error
	return &order, err
}

func (r *pharmacyRepository) GetAllOrders() ([]models.PharmacyOrder, error) {
	var orders []models.PharmacyOrder
	err := r.db.Preload("Items").Preload("Items.Medicine").Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *pharmacyRepository) CreateOrder(order *models.PharmacyOrder) error {
	return r.db.Create(order).Error
}

func (r *pharmacyRepository) UpdateOrderStatus(id uint, status string) error {
	return r.db.Model(&models.PharmacyOrder{}).Where("id = ?", id).Update("status", status).Error
}

func (r *pharmacyRepository) UpdateOrderPayment(id uint, paymentStatus string) error {
	return r.db.Model(&models.PharmacyOrder{}).Where("id = ?", id).Update("payment_status", paymentStatus).Error
}
