package services

import (
	"backend/models"
	"backend/repositories"
	"errors"
	"fmt"
	"time"
)

type PharmacyService interface {
	// Medicine Categories
	GetAllCategories() ([]models.MedicineCategory, error)
	GetCategoryByID(id uint) (*models.MedicineCategory, error)
	CreateCategory(cat *models.MedicineCategory) error
	UpdateCategory(cat *models.MedicineCategory) error
	DeleteCategory(id uint) error

	// Medicines
	GetAllMedicines(categoryID uint, search string) ([]models.Medicine, error)
	GetMedicineByID(id uint) (*models.Medicine, error)
	CreateMedicine(med *models.Medicine) error
	UpdateMedicine(med *models.Medicine) error
	DeleteMedicine(id uint) error

	// Prescriptions
	GetPrescriptionsByUserID(userID uint) ([]models.Prescription, error)
	GetPrescriptionByID(id uint) (*models.Prescription, error)
	GetAllPrescriptions() ([]models.Prescription, error)
	CreatePrescription(p *models.Prescription) error
	UpdatePrescriptionStatus(id uint, status string) error

	// Orders
	GetOrdersByUserID(userID uint) ([]models.PharmacyOrder, error)
	GetOrderByID(id uint) (*models.PharmacyOrder, error)
	GetAllOrders() ([]models.PharmacyOrder, error)
	CreateOrder(order *models.PharmacyOrder) error
	UpdateOrderStatus(id uint, status string) error
	UpdateOrderPayment(id uint, paymentStatus string) error

	// Cart
	GetCartByUserID(userID uint) (*models.Cart, error)
	AddToCart(userID uint, medicineID uint, quantity int) error
	UpdateCartItemQuantity(userID uint, cartItemID uint, quantity int) error
	RemoveFromCart(userID uint, cartItemID uint) error
	ClearCart(userID uint) error
}

type pharmacyService struct {
	repo repositories.PharmacyRepository
}

func NewPharmacyService(repo repositories.PharmacyRepository) PharmacyService {
	return &pharmacyService{repo}
}

func (s *pharmacyService) GetAllCategories() ([]models.MedicineCategory, error) {
	return s.repo.GetAllCategories()
}

func (s *pharmacyService) GetCategoryByID(id uint) (*models.MedicineCategory, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *pharmacyService) CreateCategory(cat *models.MedicineCategory) error {
	return s.repo.CreateCategory(cat)
}

func (s *pharmacyService) UpdateCategory(cat *models.MedicineCategory) error {
	return s.repo.UpdateCategory(cat)
}

func (s *pharmacyService) DeleteCategory(id uint) error {
	return s.repo.DeleteCategory(id)
}

func (s *pharmacyService) GetAllMedicines(categoryID uint, search string) ([]models.Medicine, error) {
	return s.repo.GetAllMedicines(categoryID, search)
}

func (s *pharmacyService) GetMedicineByID(id uint) (*models.Medicine, error) {
	return s.repo.GetMedicineByID(id)
}

func (s *pharmacyService) CreateMedicine(med *models.Medicine) error {
	return s.repo.CreateMedicine(med)
}

func (s *pharmacyService) UpdateMedicine(med *models.Medicine) error {
	return s.repo.UpdateMedicine(med)
}

func (s *pharmacyService) DeleteMedicine(id uint) error {
	return s.repo.DeleteMedicine(id)
}

func (s *pharmacyService) GetPrescriptionsByUserID(userID uint) ([]models.Prescription, error) {
	return s.repo.GetPrescriptionsByUserID(userID)
}

func (s *pharmacyService) GetPrescriptionByID(id uint) (*models.Prescription, error) {
	return s.repo.GetPrescriptionByID(id)
}

func (s *pharmacyService) GetAllPrescriptions() ([]models.Prescription, error) {
	return s.repo.GetAllPrescriptions()
}

func (s *pharmacyService) CreatePrescription(p *models.Prescription) error {
	p.Status = "Active"
	return s.repo.CreatePrescription(p)
}

func (s *pharmacyService) UpdatePrescriptionStatus(id uint, status string) error {
	return s.repo.UpdatePrescriptionStatus(id, status)
}

func (s *pharmacyService) GetOrdersByUserID(userID uint) ([]models.PharmacyOrder, error) {
	return s.repo.GetOrdersByUserID(userID)
}

func (s *pharmacyService) GetOrderByID(id uint) (*models.PharmacyOrder, error) {
	return s.repo.GetOrderByID(id)
}

func (s *pharmacyService) GetAllOrders() ([]models.PharmacyOrder, error) {
	return s.repo.GetAllOrders()
}

func (s *pharmacyService) CreateOrder(order *models.PharmacyOrder) error {
	if len(order.Items) == 0 {
		return errors.New("order must have at least one item")
	}

	// Generate Order Number
	order.OrderNumber = fmt.Sprintf("ORD-%d", time.Now().Unix())
	order.Status = "Pending"
	order.PaymentStatus = "Unpaid"

	// Validate items, calculate total, and update stock
	var total float64
	for i, item := range order.Items {
		med, err := s.repo.GetMedicineByID(item.MedicineID)
		if err != nil {
			return fmt.Errorf("medicine %d not found", item.MedicineID)
		}
		if med.Stock < item.Quantity {
			return fmt.Errorf("insufficient stock for %s", med.Name)
		}
		
		// Set price and subtotal
		order.Items[i].Price = med.Price
		order.Items[i].Subtotal = med.Price * float64(item.Quantity)
		total += order.Items[i].Subtotal
	}

	order.TotalPrice = total

	// Create order in DB
	if err := s.repo.CreateOrder(order); err != nil {
		return err
	}

	// Deduct stock
	for _, item := range order.Items {
		s.repo.UpdateMedicineStock(item.MedicineID, -item.Quantity)
	}

	// Update prescription status if redeemed
	if order.PrescriptionID != nil {
		s.repo.UpdatePrescriptionStatus(*order.PrescriptionID, "Redeemed")
	}

	return nil
}

func (s *pharmacyService) UpdateOrderStatus(id uint, status string) error {
	return s.repo.UpdateOrderStatus(id, status)
}

func (s *pharmacyService) UpdateOrderPayment(id uint, paymentStatus string) error {
	return s.repo.UpdateOrderPayment(id, paymentStatus)
}

func (s *pharmacyService) GetCartByUserID(userID uint) (*models.Cart, error) {
	return s.repo.GetCartByUserID(userID)
}

func (s *pharmacyService) AddToCart(userID uint, medicineID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	return s.repo.AddToCart(userID, medicineID, quantity)
}

func (s *pharmacyService) UpdateCartItemQuantity(userID uint, cartItemID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	return s.repo.UpdateCartItemQuantity(userID, cartItemID, quantity)
}

func (s *pharmacyService) RemoveFromCart(userID uint, cartItemID uint) error {
	return s.repo.RemoveFromCart(userID, cartItemID)
}

func (s *pharmacyService) ClearCart(userID uint) error {
	return s.repo.ClearCart(userID)
}
