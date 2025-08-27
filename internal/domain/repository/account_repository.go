package repository

import (
	"context"

	"savvy-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
}
