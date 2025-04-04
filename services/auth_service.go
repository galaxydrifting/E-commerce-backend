package services

import (
	"errors"
	"os"
	"time"

	"e-commerce/models"
	"e-commerce/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(name, email, password string) error {
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	// 在創建用戶前先進行密碼雜湊
	if err := user.HashPassword(); err != nil {
		return err
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := user.ComparePassword(password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, "", errors.New("could not generate token")
	}

	return user, tokenString, nil
}
