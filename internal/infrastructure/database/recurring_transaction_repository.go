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
)

type recurringTransactionRepository struct {
	db *sql.DB
}

func NewRecurringTransactionRepository(db *sql.DB) repository.RecurringTransactionRepository {
	return &recurringTransactionRepository{db: db}
}

func (r *recurringTransactionRepository) Create(ctx context.Context, recurring *entity.RecurringTransaction) error {
	query := `
		INSERT INTO recurring_transactions (
			id, user_id, category_id, account_id, amount, type, note, frequency,
			start_date, end_date, next_execution_date, last_execution_date,
			is_active, auto_execute, remaining_executions, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := r.db.ExecContext(ctx, query,
		recurring.ID,
		recurring.UserID,
		recurring.CategoryID,
		recurring.AccountID,
		recurring.Amount,
		recurring.Type,
		recurring.Note,
		recurring.Frequency,
		recurring.StartDate,
		recurring.EndDate,
		recurring.NextExecutionDate,
		recurring.LastExecutionDate,
		recurring.IsActive,
		recurring.AutoExecute,
		recurring.RemainingExecutions,
		recurring.CreatedAt,
		recurring.UpdatedAt,
	)

	return err
}

func (r *recurringTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.RecurringTransaction, error) {
	query := `
		SELECT id, user_id, category_id, account_id, amount, type, note, frequency,
			   start_date, end_date, next_execution_date, last_execution_date,
			   is_active, auto_execute, remaining_executions, created_at, updated_at
		FROM recurring_transactions WHERE id = $1
	`

	recurring := &entity.RecurringTransaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&recurring.ID,
		&recurring.UserID,
		&recurring.CategoryID,
		&recurring.AccountID,
		&recurring.Amount,
		&recurring.Type,
		&recurring.Note,
		&recurring.Frequency,
		&recurring.StartDate,
		&recurring.EndDate,
		&recurring.NextExecutionDate,
		&recurring.LastExecutionDate,
		&recurring.IsActive,
		&recurring.AutoExecute,
		&recurring.RemainingExecutions,
		&recurring.CreatedAt,
		&recurring.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return recurring, nil
}

func (r *recurringTransactionRepository) GetByFilter(ctx context.Context, filter repository.RecurringTransactionFilter) ([]*entity.RecurringTransaction, error) {
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

	if filter.AccountID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("account_id = $%d", argCount))
		args = append(args, *filter.AccountID)
	}

	if filter.Type != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("type = $%d", argCount))
		args = append(args, *filter.Type)
	}

	if filter.Frequency != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("frequency = $%d", argCount))
		args = append(args, *filter.Frequency)
	}

	if filter.IsActive != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argCount))
		args = append(args, *filter.IsActive)
	}

	if filter.DueDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("next_execution_date <= $%d", argCount))
		args = append(args, *filter.DueDate)
	}

	query := `
		SELECT id, user_id, category_id, account_id, amount, type, note, frequency,
			   start_date, end_date, next_execution_date, last_execution_date,
			   is_active, auto_execute, remaining_executions, created_at, updated_at
		FROM recurring_transactions 
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY next_execution_date ASC, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entity.RecurringTransaction
	for rows.Next() {
		recurring := &entity.RecurringTransaction{}
		err := rows.Scan(
			&recurring.ID,
			&recurring.UserID,
			&recurring.CategoryID,
			&recurring.AccountID,
			&recurring.Amount,
			&recurring.Type,
			&recurring.Note,
			&recurring.Frequency,
			&recurring.StartDate,
			&recurring.EndDate,
			&recurring.NextExecutionDate,
			&recurring.LastExecutionDate,
			&recurring.IsActive,
			&recurring.AutoExecute,
			&recurring.RemainingExecutions,
			&recurring.CreatedAt,
			&recurring.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, recurring)
	}

	return transactions, nil
}

func (r *recurringTransactionRepository) GetDueTransactions(ctx context.Context, date time.Time) ([]*entity.RecurringTransaction, error) {
	query := `
		SELECT id, user_id, category_id, account_id, amount, type, note, frequency,
			   start_date, end_date, next_execution_date, last_execution_date,
			   is_active, auto_execute, remaining_executions, created_at, updated_at
		FROM recurring_transactions 
		WHERE is_active = true 
		  AND next_execution_date <= $1
		  AND (end_date IS NULL OR end_date >= $1)
		  AND (remaining_executions IS NULL OR remaining_executions > 0)
		ORDER BY next_execution_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entity.RecurringTransaction
	for rows.Next() {
		recurring := &entity.RecurringTransaction{}
		err := rows.Scan(
			&recurring.ID,
			&recurring.UserID,
			&recurring.CategoryID,
			&recurring.AccountID,
			&recurring.Amount,
			&recurring.Type,
			&recurring.Note,
			&recurring.Frequency,
			&recurring.StartDate,
			&recurring.EndDate,
			&recurring.NextExecutionDate,
			&recurring.LastExecutionDate,
			&recurring.IsActive,
			&recurring.AutoExecute,
			&recurring.RemainingExecutions,
			&recurring.CreatedAt,
			&recurring.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, recurring)
	}

	return transactions, nil
}

func (r *recurringTransactionRepository) Update(ctx context.Context, recurring *entity.RecurringTransaction) error {
	query := `
		UPDATE recurring_transactions 
		SET category_id = $2, account_id = $3, amount = $4, type = $5, note = $6,
		    frequency = $7, start_date = $8, end_date = $9, next_execution_date = $10,
		    is_active = $11, auto_execute = $12, remaining_executions = $13, updated_at = $14
		WHERE id = $1
	`

	recurring.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		recurring.ID,
		recurring.CategoryID,
		recurring.AccountID,
		recurring.Amount,
		recurring.Type,
		recurring.Note,
		recurring.Frequency,
		recurring.StartDate,
		recurring.EndDate,
		recurring.NextExecutionDate,
		recurring.IsActive,
		recurring.AutoExecute,
		recurring.RemainingExecutions,
		recurring.UpdatedAt,
	)

	return err
}

func (r *recurringTransactionRepository) UpdateNextExecutionDate(ctx context.Context, id uuid.UUID, nextDate time.Time) error {
	query := `
		UPDATE recurring_transactions 
		SET next_execution_date = $2, updated_at = $3
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, nextDate, time.Now())
	return err
}

func (r *recurringTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM recurring_transactions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *recurringTransactionRepository) MarkAsExecuted(ctx context.Context, id uuid.UUID, executedAt time.Time) error {
	query := `
		UPDATE recurring_transactions 
		SET last_execution_date = $2, 
		    remaining_executions = CASE 
		        WHEN remaining_executions IS NOT NULL THEN remaining_executions - 1 
		        ELSE NULL 
		    END,
		    updated_at = $3
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, executedAt, time.Now())
	return err
}
