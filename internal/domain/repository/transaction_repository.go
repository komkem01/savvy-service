package repository

import (
	"context"
	"time"

	"savvy-backend/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionFilter struct {
	UserID      uuid.UUID
	AccountID   *uuid.UUID
	CategoryID  *uuid.UUID
	Type        *entity.TransactionType
	StartDate   *time.Time
	EndDate     *time.Time
	SearchQuery *string // ค้นหาจาก note หรือชื่อหมวดหมู่
	MinAmount   *decimal.Decimal
	MaxAmount   *decimal.Decimal
	Limit       int
	Offset      int
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetByFilter(ctx context.Context, filter TransactionFilter) ([]*entity.Transaction, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetMonthlySpending(ctx context.Context, userID uuid.UUID, year int, month int) (map[uuid.UUID]float64, error)
}
