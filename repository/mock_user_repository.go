package repository

import (
	"e-commerce/models"
	"errors"
)

type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() UserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(user *models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
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
	// Find user by ID first
	var found bool
	for _, existingUser := range m.users {
		if existingUser.ID == user.ID {
			found = true
			break
		}
	}
	if !found {
		return errors.New("user not found")
	}

	// If email changed, check for duplicates
	if user.Email != "" {
		for _, existingUser := range m.users {
			if existingUser.ID != user.ID && existingUser.Email == user.Email {
				return errors.New("email already in use")
			}
		}
	}

	m.users[user.Email] = user
	return nil
}
