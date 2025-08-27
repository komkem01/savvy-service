package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BudgetPeriod string

const (
	BudgetPeriodMonthly BudgetPeriod = "monthly"
	BudgetPeriodYearly  BudgetPeriod = "yearly"
)

type Budget struct {
	ID         uuid.UUID       `json:"id"`
	UserID     uuid.UUID       `json:"user_id"`
	CategoryID uuid.UUID       `json:"category_id"`
	Amount     decimal.Decimal `json:"amount"`
	Period     BudgetPeriod    `json:"period"`
	StartDate  time.Time       `json:"start_date"`
	EndDate    *time.Time      `json:"end_date,omitempty"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type BudgetProgress struct {
	BudgetID           uuid.UUID       `json:"budget_id"`
	Budget             *Budget         `json:"budget"`
	CategoryName       string          `json:"category_name"`
	BudgetAmount       decimal.Decimal `json:"budget_amount"`
	SpentAmount        decimal.Decimal `json:"spent_amount"`
	RemainingAmount    decimal.Decimal `json:"remaining_amount"`
	ProgressPercentage float64         `json:"progress_percentage"`
	IsOverBudget       bool            `json:"is_over_budget"`
	Period             string          `json:"period"` // e.g., "2024-08" for August 2024
}

func NewBudget(userID, categoryID uuid.UUID, amount decimal.Decimal, period BudgetPeriod, startDate time.Time) *Budget {
	return &Budget{
		ID:         uuid.New(),
		UserID:     userID,
		CategoryID: categoryID,
		Amount:     amount,
		Period:     period,
		StartDate:  startDate,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
