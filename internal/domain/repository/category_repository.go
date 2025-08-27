package repository

import (
	"context"

	"savvy-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type CategoryFilter struct {
	UserID     *uuid.UUID
	IsSystem   *bool
	IsArchived *bool
	Type       *entity.TransactionType
	SearchName *string
}

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error)
	GetSystemCategories(ctx context.Context) ([]*entity.Category, error)
	GetByFilter(ctx context.Context, filter CategoryFilter) ([]*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Archive(ctx context.Context, id uuid.UUID) error
	Unarchive(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetCategoryUsageStats(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]int64, error)
}
