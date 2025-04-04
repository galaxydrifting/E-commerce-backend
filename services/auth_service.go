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

func (s *AuthService) Logout(userID uint) error {
	// 由於這是 demo 專案，這裡不實現 token 撤銷
	// 在實際專案中，可以：
	// 1. 使用 Redis 維護一個 token 黑名單
	// 2. 使用較短的 token 過期時間，並實現 refresh token 機制
	// 3. 在資料庫中維護 token 狀態
	// 4. 使用第三方服務如 Auth0 來處理完整的身份驗證流程
	return nil
}

func (s *AuthService) UpdateProfile(userID uint, name, email string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if email is already taken by another user
	if user.Email != email {
		existingUser, _ := s.userRepo.FindByEmail(email)
		if existingUser != nil {
			return nil, errors.New("email already in use")
		}
	}

	user.Name = name
	user.Email = email

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, errors.New("failed to update profile")
	}

	return user, nil
}
