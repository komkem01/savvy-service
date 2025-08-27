package repository

import (
	"context"
	"time"

	"savvy-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type RecurringTransactionFilter struct {
	UserID     uuid.UUID
	CategoryID *uuid.UUID
	AccountID  *uuid.UUID
	Type       *entity.TransactionType
	Frequency  *entity.RecurringFrequency
	IsActive   *bool
	DueDate    *time.Time // สำหรับหารายการที่ถึงเวลาต้องรัน
}

type RecurringTransactionRepository interface {
	Create(ctx context.Context, recurring *entity.RecurringTransaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.RecurringTransaction, error)
	GetByFilter(ctx context.Context, filter RecurringTransactionFilter) ([]*entity.RecurringTransaction, error)
	GetDueTransactions(ctx context.Context, date time.Time) ([]*entity.RecurringTransaction, error)
	Update(ctx context.Context, recurring *entity.RecurringTransaction) error
	UpdateNextExecutionDate(ctx context.Context, id uuid.UUID, nextDate time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
	MarkAsExecuted(ctx context.Context, id uuid.UUID, executedAt time.Time) error
}
