package services

import (
	"backend/models"
	"errors"
)

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
