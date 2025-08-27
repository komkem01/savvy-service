package repository

import (
	"context"

	"savvy-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type SavingsGoalRepository interface {
	Create(ctx context.Context, goal *entity.SavingsGoal) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SavingsGoal, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SavingsGoal, error)
	Update(ctx context.Context, goal *entity.SavingsGoal) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GoalDepositRepository interface {
	Create(ctx context.Context, deposit *entity.GoalDeposit) error
	GetByGoalID(ctx context.Context, goalID uuid.UUID) ([]*entity.GoalDeposit, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.GoalDeposit, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
