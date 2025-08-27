package repository

import (
	"context"
	"time"

	"savvy-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID, loginTime time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
}
