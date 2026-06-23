package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

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
