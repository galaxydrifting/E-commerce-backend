package services

import (
	"e-commerce/models"
	"errors"
	"testing"
	"time"
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

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Update(user *models.User) error {
	// 檢查是否存在具有相同 ID 的用戶
	found := false
	for _, existingUser := range m.users {
		if existingUser.ID == user.ID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("user not found")
	}

	// 更新用戶資料
	oldEmail := ""
	for email, existingUser := range m.users {
		if existingUser.ID == user.ID {
			oldEmail = email
			break
		}
	}

	// 如果電子郵件已更改，確保新的電子郵件未被其他用戶使用
	if oldEmail != user.Email {
		for _, existingUser := range m.users {
			if existingUser.ID != user.ID && existingUser.Email == user.Email {
				return errors.New("email already in use")
			}
		}
		delete(m.users, oldEmail)
	}

	m.users[user.Email] = user
	return nil
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

func TestUpdateProfile(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create test users
	testUser := &models.User{
		ID:        1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.users[testUser.Email] = testUser

	anotherUser := &models.User{
		ID:        2,
		Name:     "Another User",
		Email:    "another@example.com",
		Password: "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.users[anotherUser.Email] = anotherUser

	tests := []struct {
		name     string
		userID   uint
		newName  string
		newEmail string
		wantErr  bool
	}{
		{
			name:     "successful update",
			userID:   1,
			newName:  "Updated Name",
			newEmail: "updated@example.com",
			wantErr:  false,
		},
		{
			name:     "non-existent user",
			userID:   99,
			newName:  "New Name",
			newEmail: "new@example.com",
			wantErr:  true,
		},
		{
			name:     "email already taken",
			userID:   1,
			newName:  "Test User",
			newEmail: "another@example.com",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedUser, err := authService.UpdateProfile(tt.userID, tt.newName, tt.newEmail)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if updatedUser == nil {
					t.Error("UpdateProfile() returned nil user for successful update")
					return
				}
				if updatedUser.Name != tt.newName {
					t.Errorf("UpdateProfile() name = %v, want %v", updatedUser.Name, tt.newName)
				}
				if updatedUser.Email != tt.newEmail {
					t.Errorf("UpdateProfile() email = %v, want %v", updatedUser.Email, tt.newEmail)
				}
			}
		})
	}
}
