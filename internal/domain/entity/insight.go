package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InsightType string

const (
	InsightTypeSpending               InsightType = "spending"
	InsightTypeBudget                 InsightType = "budget"
	InsightTypeGoal                   InsightType = "goal"
	InsightTypeTrend                  InsightType = "trend"
	InsightTypeRecommend              InsightType = "recommendation"
	InsightTypeAnomalyDetection       InsightType = "anomaly_detection"
	InsightTypeSpendingPattern        InsightType = "spending_pattern"
	InsightTypeBudgetAlert            InsightType = "budget_alert"
	InsightTypeSavingsRecommendation  InsightType = "savings_recommendation"
	InsightTypeCategoryRecommendation InsightType = "category_recommendation"
)

type InsightPriority string

const (
	InsightPriorityLow    InsightPriority = "low"
	InsightPriorityMedium InsightPriority = "medium"
	InsightPriorityHigh   InsightPriority = "high"
)

type Insight struct {
	ID                uuid.UUID       `json:"id" db:"id"`
	UserID            uuid.UUID       `json:"user_id" db:"user_id"`
	Type              InsightType     `json:"type" db:"type"`
	Priority          InsightPriority `json:"priority" db:"priority"`
	Title             string          `json:"title" db:"title"`
	Content           string          `json:"content" db:"content"`
	ActionText        *string         `json:"action_text,omitempty" db:"action_text"`
	IsRead            bool            `json:"is_read" db:"is_read"`
	RelatedEntityID   *uuid.UUID      `json:"related_entity_id,omitempty" db:"related_entity_id"`
	RelatedEntityType *string         `json:"related_entity_type,omitempty" db:"related_entity_type"`
	RelatedData       json.RawMessage `json:"related_data,omitempty" db:"related_data"`
	ValidUntil        *time.Time      `json:"valid_until,omitempty" db:"valid_until"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
}

type SpendingAnomaly struct {
	CategoryID         uuid.UUID       `json:"category_id"`
	CategoryName       string          `json:"category_name"`
	CurrentAmount      decimal.Decimal `json:"current_amount"`
	AverageAmount      decimal.Decimal `json:"average_amount"`
	PercentageIncrease float64         `json:"percentage_increase"`
	Period             string          `json:"period"`
	Severity           string          `json:"severity"` // "low", "medium", "high"
}

type SpendingPattern struct {
	CategoryID     uuid.UUID       `json:"category_id"`
	CategoryName   string          `json:"category_name"`
	DayOfWeek      string          `json:"day_of_week"`
	TimeOfDay      string          `json:"time_of_day"`
	FrequencyCount int             `json:"frequency_count"`
	AverageAmount  decimal.Decimal `json:"average_amount"`
	Period         string          `json:"period"`
}

func NewInsight(userID uuid.UUID, insightType InsightType, content string, relatedData json.RawMessage) *Insight {
	return &Insight{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        insightType,
		Priority:    InsightPriorityMedium,
		Title:       "",
		Content:     content,
		IsRead:      false,
		RelatedData: relatedData,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewAdvancedInsight(userID uuid.UUID, insightType InsightType, priority InsightPriority, title, content string) *Insight {
	return &Insight{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      insightType,
		Priority:  priority,
		Title:     title,
		Content:   content,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
