package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	UserID          uuid.UUID       `json:"user_id" db:"user_id"`
	CategoryID      uuid.UUID       `json:"category_id" db:"category_id"`
	AccountID       uuid.UUID       `json:"account_id" db:"account_id"`
	Amount          decimal.Decimal `json:"amount" db:"amount"`
	Type            TransactionType `json:"type" db:"type"`
	Note            *string         `json:"note,omitempty" db:"note"`
	TransactionDate time.Time       `json:"transaction_date" db:"transaction_date"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

func NewTransaction(userID, categoryID, accountID uuid.UUID, amount decimal.Decimal,
	transactionType TransactionType, note *string, transactionDate time.Time) *Transaction {
	return &Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		CategoryID:      categoryID,
		AccountID:       accountID,
		Amount:          amount,
		Type:            transactionType,
		Note:            note,
		TransactionDate: transactionDate,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
