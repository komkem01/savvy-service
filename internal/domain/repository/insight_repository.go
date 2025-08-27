package repository

import (
	"context"
	"time"

	"savvy-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type InsightFilter struct {
	UserID    uuid.UUID
	Type      *entity.InsightType
	Priority  *entity.InsightPriority
	IsRead    *bool
	ValidOnly bool // filter out expired insights
	Limit     int
	Offset    int
}

type InsightRepository interface {
	Create(ctx context.Context, insight *entity.Insight) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Insight, error)
	GetByFilter(ctx context.Context, filter InsightFilter) ([]*entity.Insight, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.Insight, error)
	GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Insight, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context, before time.Time) error
	GetSpendingAnomalies(ctx context.Context, userID uuid.UUID, months int) ([]*entity.SpendingAnomaly, error)
	GetSpendingPatterns(ctx context.Context, userID uuid.UUID, days int) ([]*entity.SpendingPattern, error)
}
