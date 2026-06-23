package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
	AddFamilyMember(member *models.FamilyMember) error
	DeleteFamilyMember(id uint) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("FamilyMembers").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *authRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("FamilyMembers").First(&user, id).Error
	return &user, err
}

func (r *authRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *authRepository) AddFamilyMember(member *models.FamilyMember) error {
	return r.db.Create(member).Error
}

func (r *authRepository) DeleteFamilyMember(id uint) error {
	return r.db.Delete(&models.FamilyMember{}, id).Error
}
