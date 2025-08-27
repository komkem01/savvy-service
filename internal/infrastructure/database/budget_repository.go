package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type budgetRepository struct {
	db *sql.DB
}

func NewBudgetRepository(db *sql.DB) repository.BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) Create(ctx context.Context, budget *entity.Budget) error {
	query := `
		INSERT INTO budgets (id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		budget.ID,
		budget.UserID,
		budget.CategoryID,
		budget.Amount,
		budget.Period,
		budget.StartDate,
		budget.EndDate,
		budget.IsActive,
		budget.CreatedAt,
		budget.UpdatedAt,
	)

	return err
}

func (r *budgetRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Budget, error) {
	query := `
		SELECT id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
		FROM budgets WHERE id = $1
	`

	budget := &entity.Budget{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&budget.ID,
		&budget.UserID,
		&budget.CategoryID,
		&budget.Amount,
		&budget.Period,
		&budget.StartDate,
		&budget.EndDate,
		&budget.IsActive,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return budget, nil
}

func (r *budgetRepository) GetByFilter(ctx context.Context, filter repository.BudgetFilter) ([]*entity.Budget, error) {
	var conditions []string
	var args []interface{}
	argCount := 0

	// Build WHERE clause
	conditions = append(conditions, "user_id = $1")
	args = append(args, filter.UserID)
	argCount = 1

	if filter.CategoryID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argCount))
		args = append(args, *filter.CategoryID)
	}

	if filter.Period != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("period = $%d", argCount))
		args = append(args, *filter.Period)
	}

	if filter.IsActive != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argCount))
		args = append(args, *filter.IsActive)
	}

	if filter.StartDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", argCount))
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("(end_date IS NULL OR end_date <= $%d)", argCount))
		args = append(args, *filter.EndDate)
	}

	query := `
		SELECT id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
		FROM budgets 
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []*entity.Budget
	for rows.Next() {
		budget := &entity.Budget{}
		err := rows.Scan(
			&budget.ID,
			&budget.UserID,
			&budget.CategoryID,
			&budget.Amount,
			&budget.Period,
			&budget.StartDate,
			&budget.EndDate,
			&budget.IsActive,
			&budget.CreatedAt,
			&budget.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

func (r *budgetRepository) GetByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID) (*entity.Budget, error) {
	query := `
		SELECT id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at
		FROM budgets 
		WHERE user_id = $1 AND category_id = $2 AND is_active = true
		ORDER BY created_at DESC
		LIMIT 1
	`

	budget := &entity.Budget{}
	err := r.db.QueryRowContext(ctx, query, userID, categoryID).Scan(
		&budget.ID,
		&budget.UserID,
		&budget.CategoryID,
		&budget.Amount,
		&budget.Period,
		&budget.StartDate,
		&budget.EndDate,
		&budget.IsActive,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return budget, nil
}

func (r *budgetRepository) Update(ctx context.Context, budget *entity.Budget) error {
	query := `
		UPDATE budgets 
		SET amount = $2, period = $3, start_date = $4, end_date = $5, is_active = $6, updated_at = $7
		WHERE id = $1
	`

	budget.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		budget.ID,
		budget.Amount,
		budget.Period,
		budget.StartDate,
		budget.EndDate,
		budget.IsActive,
		budget.UpdatedAt,
	)

	return err
}

func (r *budgetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM budgets WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *budgetRepository) GetBudgetProgress(ctx context.Context, userID uuid.UUID, year, month int) ([]*entity.BudgetProgress, error) {
	query := `
		SELECT 
			b.id as budget_id,
			b.amount as budget_amount,
			c.name as category_name,
			COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0) as spent_amount
		FROM budgets b
		INNER JOIN categories c ON b.category_id = c.id
		LEFT JOIN transactions t ON b.category_id = t.category_id 
			AND t.user_id = b.user_id 
			AND EXTRACT(YEAR FROM t.transaction_date) = $2
			AND EXTRACT(MONTH FROM t.transaction_date) = $3
		WHERE b.user_id = $1 AND b.is_active = true
		GROUP BY b.id, b.amount, c.name
		ORDER BY c.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.BudgetProgress
	for rows.Next() {
		var budgetID uuid.UUID
		var budgetAmount, spentAmount decimal.Decimal
		var categoryName string

		err := rows.Scan(&budgetID, &budgetAmount, &categoryName, &spentAmount)
		if err != nil {
			return nil, err
		}

		remainingAmount := budgetAmount.Sub(spentAmount)
		progressPercentage := 0.0
		if !budgetAmount.IsZero() {
			progressPercentage, _ = spentAmount.Div(budgetAmount).Mul(decimal.NewFromInt(100)).Float64()
		}

		progress := &entity.BudgetProgress{
			BudgetID:           budgetID,
			CategoryName:       categoryName,
			BudgetAmount:       budgetAmount,
			SpentAmount:        spentAmount,
			RemainingAmount:    remainingAmount,
			ProgressPercentage: progressPercentage,
			IsOverBudget:       spentAmount.GreaterThan(budgetAmount),
			Period:             fmt.Sprintf("%d-%02d", year, month),
		}

		results = append(results, progress)
	}

	return results, nil
}

func (r *budgetRepository) GetBudgetProgressByCategory(ctx context.Context, userID, categoryID uuid.UUID, year, month int) (*entity.BudgetProgress, error) {
	query := `
		SELECT 
			b.id as budget_id,
			b.amount as budget_amount,
			c.name as category_name,
			COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0) as spent_amount
		FROM budgets b
		INNER JOIN categories c ON b.category_id = c.id
		LEFT JOIN transactions t ON b.category_id = t.category_id 
			AND t.user_id = b.user_id 
			AND EXTRACT(YEAR FROM t.transaction_date) = $3
			AND EXTRACT(MONTH FROM t.transaction_date) = $4
		WHERE b.user_id = $1 AND b.category_id = $2 AND b.is_active = true
		GROUP BY b.id, b.amount, c.name
		LIMIT 1
	`

	var budgetID uuid.UUID
	var budgetAmount, spentAmount decimal.Decimal
	var categoryName string

	err := r.db.QueryRowContext(ctx, query, userID, categoryID, year, month).Scan(
		&budgetID, &budgetAmount, &categoryName, &spentAmount,
	)
	if err != nil {
		return nil, err
	}

	remainingAmount := budgetAmount.Sub(spentAmount)
	progressPercentage := 0.0
	if !budgetAmount.IsZero() {
		progressPercentage, _ = spentAmount.Div(budgetAmount).Mul(decimal.NewFromInt(100)).Float64()
	}

	return &entity.BudgetProgress{
		BudgetID:           budgetID,
		CategoryName:       categoryName,
		BudgetAmount:       budgetAmount,
		SpentAmount:        spentAmount,
		RemainingAmount:    remainingAmount,
		ProgressPercentage: progressPercentage,
		IsOverBudget:       spentAmount.GreaterThan(budgetAmount),
		Period:             fmt.Sprintf("%d-%02d", year, month),
	}, nil
}
