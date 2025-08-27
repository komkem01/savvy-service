package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountType string

const (
	AccountTypeCash    AccountType = "cash"
	AccountTypeBank    AccountType = "bank"
	AccountTypeCredit  AccountType = "credit"
	AccountTypeSavings AccountType = "savings"
)

type Account struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	UserID         uuid.UUID       `json:"user_id" db:"user_id"`
	Name           string          `json:"name" db:"name"`
	Type           AccountType     `json:"type" db:"type"`
	InitialBalance decimal.Decimal `json:"initial_balance" db:"initial_balance"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

func NewAccount(userID uuid.UUID, name string, accountType AccountType, initialBalance decimal.Decimal) *Account {
	return &Account{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           name,
		Type:           accountType,
		InitialBalance: initialBalance,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
