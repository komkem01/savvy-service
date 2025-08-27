package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RecurringFrequency string

const (
	RecurringFrequencyDaily   RecurringFrequency = "daily"
	RecurringFrequencyWeekly  RecurringFrequency = "weekly"
	RecurringFrequencyMonthly RecurringFrequency = "monthly"
	RecurringFrequencyYearly  RecurringFrequency = "yearly"
)

type RecurringTransaction struct {
	ID                  uuid.UUID          `json:"id"`
	UserID              uuid.UUID          `json:"user_id"`
	CategoryID          uuid.UUID          `json:"category_id"`
	AccountID           uuid.UUID          `json:"account_id"`
	Amount              decimal.Decimal    `json:"amount"`
	Type                TransactionType    `json:"type"`
	Note                *string            `json:"note,omitempty"`
	Frequency           RecurringFrequency `json:"frequency"`
	StartDate           time.Time          `json:"start_date"`
	EndDate             *time.Time         `json:"end_date,omitempty"`
	NextExecutionDate   time.Time          `json:"next_execution_date"`
	LastExecutionDate   *time.Time         `json:"last_execution_date,omitempty"`
	IsActive            bool               `json:"is_active"`
	AutoExecute         bool               `json:"auto_execute"`
	RemainingExecutions *int               `json:"remaining_executions,omitempty"` // null = unlimited
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
}

func NewRecurringTransaction(
	userID, categoryID, accountID uuid.UUID,
	amount decimal.Decimal,
	transactionType TransactionType,
	frequency RecurringFrequency,
	startDate time.Time,
	note *string,
) *RecurringTransaction {
	return &RecurringTransaction{
		ID:                uuid.New(),
		UserID:            userID,
		CategoryID:        categoryID,
		AccountID:         accountID,
		Amount:            amount,
		Type:              transactionType,
		Note:              note,
		Frequency:         frequency,
		StartDate:         startDate,
		NextExecutionDate: startDate,
		IsActive:          true,
		AutoExecute:       false, // Default to manual approval
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// CalculateNextExecutionDate คำนวณวันที่รันครั้งถัดไป
func (rt *RecurringTransaction) CalculateNextExecutionDate() time.Time {
	baseDate := rt.NextExecutionDate
	if rt.LastExecutionDate != nil {
		baseDate = *rt.LastExecutionDate
	}

	switch rt.Frequency {
	case RecurringFrequencyDaily:
		return baseDate.AddDate(0, 0, 1)
	case RecurringFrequencyWeekly:
		return baseDate.AddDate(0, 0, 7)
	case RecurringFrequencyMonthly:
		return baseDate.AddDate(0, 1, 0)
	case RecurringFrequencyYearly:
		return baseDate.AddDate(1, 0, 0)
	default:
		return baseDate.AddDate(0, 1, 0) // Default to monthly
	}
}
