package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	Email              string     `json:"email" db:"email"`
	PasswordHash       string     `json:"-" db:"password_hash"` // Don't include in JSON
	DisplayName        *string    `json:"display_name,omitempty" db:"display_name"`
	CurrencyPreference string     `json:"currency_preference" db:"currency_preference"`
	IsActive           bool       `json:"is_active" db:"is_active"`
	LastLoginAt        *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

func NewUser(email, passwordHash string, displayName *string) *User {
	return &User{
		ID:                 uuid.New(),
		Email:              email,
		PasswordHash:       passwordHash,
		DisplayName:        displayName,
		CurrencyPreference: "THB",
		IsActive:           true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
