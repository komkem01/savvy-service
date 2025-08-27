package usecase

import (
	"context"
	"fmt"
	"time"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BudgetUsecase interface {
	CreateBudget(ctx context.Context, userID, categoryID uuid.UUID, amount decimal.Decimal, period entity.BudgetPeriod, startDate time.Time, endDate *time.Time) (*entity.Budget, error)
	GetUserBudgets(ctx context.Context, userID uuid.UUID) ([]*entity.Budget, error)
	GetBudgetByID(ctx context.Context, userID, budgetID uuid.UUID) (*entity.Budget, error)
	UpdateBudget(ctx context.Context, userID uuid.UUID, budget *entity.Budget) error
	DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error
	GetBudgetProgress(ctx context.Context, userID uuid.UUID, year, month int) ([]*entity.BudgetProgress, error)
	GetCurrentMonthBudgetProgress(ctx context.Context, userID uuid.UUID) ([]*entity.BudgetProgress, error)
	CheckBudgetAlerts(ctx context.Context, userID uuid.UUID) ([]*entity.Insight, error)
}

type budgetUsecase struct {
	budgetRepo   repository.BudgetRepository
	categoryRepo repository.CategoryRepository
	insightRepo  repository.InsightRepository
}

func NewBudgetUsecase(
	budgetRepo repository.BudgetRepository,
	categoryRepo repository.CategoryRepository,
	insightRepo repository.InsightRepository,
) BudgetUsecase {
	return &budgetUsecase{
		budgetRepo:   budgetRepo,
		categoryRepo: categoryRepo,
		insightRepo:  insightRepo,
	}
}

func (b *budgetUsecase) CreateBudget(ctx context.Context, userID, categoryID uuid.UUID, amount decimal.Decimal, period entity.BudgetPeriod, startDate time.Time, endDate *time.Time) (*entity.Budget, error) {
	// Validate category belongs to user or is system category
	category, err := b.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	if category.UserID != nil && *category.UserID != userID {
		return nil, fmt.Errorf("category does not belong to user")
	}

	// Check if budget already exists for this category
	existing, _ := b.budgetRepo.GetByUserIDAndCategoryID(ctx, userID, categoryID)
	if existing != nil {
		return nil, fmt.Errorf("budget already exists for this category")
	}

	budget := entity.NewBudget(userID, categoryID, amount, period, startDate)
	budget.EndDate = endDate

	err = b.budgetRepo.Create(ctx, budget)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return budget, nil
}

func (b *budgetUsecase) GetUserBudgets(ctx context.Context, userID uuid.UUID) ([]*entity.Budget, error) {
	isActive := true
	filter := repository.BudgetFilter{
		UserID:   userID,
		IsActive: &isActive,
	}

	return b.budgetRepo.GetByFilter(ctx, filter)
}

func (b *budgetUsecase) GetBudgetByID(ctx context.Context, userID, budgetID uuid.UUID) (*entity.Budget, error) {
	budget, err := b.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		return nil, err
	}

	if budget.UserID != userID {
		return nil, fmt.Errorf("budget does not belong to user")
	}

	return budget, nil
}

func (b *budgetUsecase) UpdateBudget(ctx context.Context, userID uuid.UUID, budget *entity.Budget) error {
	existing, err := b.budgetRepo.GetByID(ctx, budget.ID)
	if err != nil {
		return err
	}

	if existing.UserID != userID {
		return fmt.Errorf("budget does not belong to user")
	}

	budget.UserID = userID // Ensure user ID is preserved
	return b.budgetRepo.Update(ctx, budget)
}

func (b *budgetUsecase) DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error {
	budget, err := b.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		return err
	}

	if budget.UserID != userID {
		return fmt.Errorf("budget does not belong to user")
	}

	return b.budgetRepo.Delete(ctx, budgetID)
}

func (b *budgetUsecase) GetBudgetProgress(ctx context.Context, userID uuid.UUID, year, month int) ([]*entity.BudgetProgress, error) {
	return b.budgetRepo.GetBudgetProgress(ctx, userID, year, month)
}

func (b *budgetUsecase) GetCurrentMonthBudgetProgress(ctx context.Context, userID uuid.UUID) ([]*entity.BudgetProgress, error) {
	now := time.Now()
	return b.budgetRepo.GetBudgetProgress(ctx, userID, now.Year(), int(now.Month()))
}

func (b *budgetUsecase) CheckBudgetAlerts(ctx context.Context, userID uuid.UUID) ([]*entity.Insight, error) {
	now := time.Now()
	progresses, err := b.budgetRepo.GetBudgetProgress(ctx, userID, now.Year(), int(now.Month()))
	if err != nil {
		return nil, err
	}

	var insights []*entity.Insight

	for _, progress := range progresses {
		// Create alert for budgets at 80% usage
		if progress.ProgressPercentage >= 80 && progress.ProgressPercentage < 100 {
			title := fmt.Sprintf("งบประมาณ %s ใกล้หมดแล้ว", progress.CategoryName)

			spentAmount, _ := progress.SpentAmount.Float64()
			budgetAmount, _ := progress.BudgetAmount.Float64()

			message := fmt.Sprintf("คุณใช้งบประมาณหมวดหมู่ %s ไปแล้ว %.1f%% (%.2f จาก %.2f บาท)",
				progress.CategoryName,
				progress.ProgressPercentage,
				spentAmount,
				budgetAmount)

			insight := entity.NewAdvancedInsight(userID, entity.InsightTypeBudgetAlert, entity.InsightPriorityMedium, title, message)
			insight.RelatedEntityID = &progress.BudgetID
			insight.RelatedEntityType = &[]string{"budget"}[0]

			validUntil := time.Now().AddDate(0, 0, 7) // Valid for 7 days
			insight.ValidUntil = &validUntil

			insights = append(insights, insight)
		}

		// Create alert for over-budget categories
		if progress.IsOverBudget {
			title := fmt.Sprintf("เกินงบประมาณ %s", progress.CategoryName)

			overAmount, _ := progress.SpentAmount.Sub(progress.BudgetAmount).Float64()

			message := fmt.Sprintf("คุณใช้จ่ายหมวดหมู่ %s เกินงบประมาณแล้ว %.2f บาท (%.1f%%)",
				progress.CategoryName,
				overAmount,
				progress.ProgressPercentage)

			insight := entity.NewAdvancedInsight(userID, entity.InsightTypeBudgetAlert, entity.InsightPriorityHigh, title, message)
			insight.RelatedEntityID = &progress.BudgetID
			insight.RelatedEntityType = &[]string{"budget"}[0]

			validUntil := time.Now().AddDate(0, 0, 3) // Valid for 3 days
			insight.ValidUntil = &validUntil

			insights = append(insights, insight)
		}
	}

	// Save insights to database
	for _, insight := range insights {
		err := b.insightRepo.Create(ctx, insight)
		if err != nil {
			// Log error but don't fail the entire operation
			continue
		}
	}

	return insights, nil
}
