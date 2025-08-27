package entity

import (
	"time"

	"github.com/google/uuid"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type Category struct {
	ID         uuid.UUID    `json:"id" db:"id"`
	UserID     *uuid.UUID   `json:"user_id,omitempty" db:"user_id"` // nullable for system categories
	Name       string       `json:"name" db:"name"`
	Type       CategoryType `json:"type" db:"type"`
	IconName   *string      `json:"icon_name,omitempty" db:"icon_name"`
	ColorHex   *string      `json:"color_hex,omitempty" db:"color_hex"`
	IsArchived bool         `json:"is_archived" db:"is_archived"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at" db:"updated_at"`
}

func NewCategory(userID *uuid.UUID, name string, categoryType CategoryType) *Category {
	return &Category{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       name,
		Type:       categoryType,
		IsArchived: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
