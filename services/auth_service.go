package services

import (
	"backend/models"
	"backend/repositories"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(fullName, email, phone, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
}

type authService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(fullName, email, phone, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FullName:    fullName,
		Email:       email,
		PhoneNumber: phone,
		Password:    string(hashedPassword),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret" // fallback
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
