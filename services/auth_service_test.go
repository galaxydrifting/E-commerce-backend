package services

import (
	"e-commerce/models"
	"errors"
	"testing"
)

// Mock UserRepository
type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(user *models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	if err := user.HashPassword(); err != nil {
		return err
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func TestRegister(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	tests := []struct {
		name     string
		userName string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "successful registration",
			userName: "Test User",
			email:    "test@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "duplicate email",
			userName: "Another User",
			email:    "test@example.com",
			password: "password123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Register(tt.userName, tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create a test user
	testUser := &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	if err := testUser.HashPassword(); err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	mockRepo.users[testUser.Email] = testUser

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "invalid email",
			email:    "wrong@example.com",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, token, err := authService.Login(tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if user == nil {
					t.Error("Login() returned nil user for successful login")
				}
				if token == "" {
					t.Error("Login() returned empty token for successful login")
				}
			}
		})
	}
}
