package repository

import (
	"context"
	"time"

	"savvy-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type BudgetFilter struct {
	UserID     uuid.UUID
	CategoryID *uuid.UUID
	Period     *entity.BudgetPeriod
	IsActive   *bool
	StartDate  *time.Time
	EndDate    *time.Time
}

type BudgetRepository interface {
	Create(ctx context.Context, budget *entity.Budget) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Budget, error)
	GetByFilter(ctx context.Context, filter BudgetFilter) ([]*entity.Budget, error)
	GetByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID) (*entity.Budget, error)
	Update(ctx context.Context, budget *entity.Budget) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetBudgetProgress(ctx context.Context, userID uuid.UUID, year, month int) ([]*entity.BudgetProgress, error)
	GetBudgetProgressByCategory(ctx context.Context, userID, categoryID uuid.UUID, year, month int) (*entity.BudgetProgress, error)
}
