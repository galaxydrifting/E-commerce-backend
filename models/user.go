package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primarykey" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt time.Time `json:"deleted_at,omitempty" gorm:"index" example:"2024-01-01T00:00:00Z"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" gorm:"unique" example:"user@example.com"`
	Password  string    `json:"password,omitempty" example:"password123"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
}
