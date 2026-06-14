package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type PharmacyRepository interface {
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
	UpdateMedicineStock(id uint, quantityChange int) error

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

type pharmacyRepository struct {
	db *gorm.DB
}

func NewPharmacyRepository(db *gorm.DB) PharmacyRepository {
	return &pharmacyRepository{db}
}

func (r *pharmacyRepository) GetAllCategories() ([]models.MedicineCategory, error) {
	var cats []models.MedicineCategory
	err := r.db.Find(&cats).Error
	return cats, err
}

func (r *pharmacyRepository) GetCategoryByID(id uint) (*models.MedicineCategory, error) {
	var cat models.MedicineCategory
	err := r.db.First(&cat, id).Error
	return &cat, err
}

func (r *pharmacyRepository) CreateCategory(cat *models.MedicineCategory) error {
	return r.db.Create(cat).Error
}

func (r *pharmacyRepository) UpdateCategory(cat *models.MedicineCategory) error {
	return r.db.Save(cat).Error
}

func (r *pharmacyRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&models.MedicineCategory{}, id).Error
}

func (r *pharmacyRepository) GetAllMedicines(categoryID uint, search string) ([]models.Medicine, error) {
	var meds []models.Medicine
	query := r.db.Preload("Category")

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name LIKE ? OR generic_name LIKE ?", searchPattern, searchPattern)
	}

	err := query.Find(&meds).Error
	return meds, err
}

func (r *pharmacyRepository) GetMedicineByID(id uint) (*models.Medicine, error) {
	var med models.Medicine
	err := r.db.Preload("Category").First(&med, id).Error
	return &med, err
}

func (r *pharmacyRepository) CreateMedicine(med *models.Medicine) error {
	return r.db.Create(med).Error
}

func (r *pharmacyRepository) UpdateMedicine(med *models.Medicine) error {
	return r.db.Save(med).Error
}

func (r *pharmacyRepository) DeleteMedicine(id uint) error {
	return r.db.Delete(&models.Medicine{}, id).Error
}

func (r *pharmacyRepository) UpdateMedicineStock(id uint, quantityChange int) error {
	return r.db.Model(&models.Medicine{}).Where("id = ?", id).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantityChange)).Error
}

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

func (r *pharmacyRepository) GetCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items").Preload("Items.Medicine").Where("user_id = ?", userID).FirstOrCreate(&cart, models.Cart{UserID: userID}).Error
	return &cart, err
}

func (r *pharmacyRepository) AddToCart(userID uint, medicineID uint, quantity int) error {
	cart, err := r.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	var item models.CartItem
	err = r.db.Where("cart_id = ? AND medicine_id = ?", cart.ID, medicineID).First(&item).Error
	if err == nil {
		item.Quantity += quantity
		return r.db.Save(&item).Error
	} else if err == gorm.ErrRecordNotFound {
		newItem := models.CartItem{
			CartID:     cart.ID,
			MedicineID: medicineID,
			Quantity:   quantity,
		}
		return r.db.Create(&newItem).Error
	}
	return err
}

func (r *pharmacyRepository) UpdateCartItemQuantity(userID uint, cartItemID uint, quantity int) error {
	cart, err := r.GetCartByUserID(userID)
	if err != nil {
		return err
	}
	return r.db.Model(&models.CartItem{}).Where("id = ? AND cart_id = ?", cartItemID, cart.ID).Update("quantity", quantity).Error
}

func (r *pharmacyRepository) RemoveFromCart(userID uint, cartItemID uint) error {
	cart, err := r.GetCartByUserID(userID)
	if err != nil {
		return err
	}
	return r.db.Where("id = ? AND cart_id = ?", cartItemID, cart.ID).Delete(&models.CartItem{}).Error
}

func (r *pharmacyRepository) ClearCart(userID uint) error {
	cart, err := r.GetCartByUserID(userID)
	if err != nil {
		return err
	}
	return r.db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
}
