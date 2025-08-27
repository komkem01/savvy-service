package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GoalStatus string

const (
	GoalStatusActive    GoalStatus = "active"
	GoalStatusCompleted GoalStatus = "completed"
	GoalStatusPaused    GoalStatus = "paused"
)

type SavingsGoal struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	UserID       uuid.UUID       `json:"user_id" db:"user_id"`
	Name         string          `json:"name" db:"name"`
	TargetAmount decimal.Decimal `json:"target_amount" db:"target_amount"`
	TargetDate   *time.Time      `json:"target_date,omitempty" db:"target_date"`
	Status       GoalStatus      `json:"status" db:"status"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

func NewSavingsGoal(userID uuid.UUID, name string, targetAmount decimal.Decimal, targetDate *time.Time) *SavingsGoal {
	return &SavingsGoal{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         name,
		TargetAmount: targetAmount,
		TargetDate:   targetDate,
		Status:       GoalStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

type GoalDeposit struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	UserID      uuid.UUID       `json:"user_id" db:"user_id"`
	GoalID      uuid.UUID       `json:"goal_id" db:"goal_id"`
	AccountID   uuid.UUID       `json:"account_id" db:"account_id"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	DepositDate time.Time       `json:"deposit_date" db:"deposit_date"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}

func NewGoalDeposit(userID, goalID, accountID uuid.UUID, amount decimal.Decimal, depositDate time.Time) *GoalDeposit {
	return &GoalDeposit{
		ID:          uuid.New(),
		UserID:      userID,
		GoalID:      goalID,
		AccountID:   accountID,
		Amount:      amount,
		DepositDate: depositDate,
		CreatedAt:   time.Now(),
	}
}
