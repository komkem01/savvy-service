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

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, category_id, account_id, amount, type, note, transaction_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		transaction.ID,
		transaction.UserID,
		transaction.CategoryID,
		transaction.AccountID,
		transaction.Amount,
		transaction.Type,
		transaction.Note,
		transaction.TransactionDate,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)

	return err
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	query := `
		SELECT id, user_id, category_id, account_id, amount, type, note, transaction_date, created_at, updated_at
		FROM transactions WHERE id = $1
	`

	transaction := &entity.Transaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.CategoryID,
		&transaction.AccountID,
		&transaction.Amount,
		&transaction.Type,
		&transaction.Note,
		&transaction.TransactionDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *transactionRepository) GetByFilter(ctx context.Context, filter repository.TransactionFilter) ([]*entity.Transaction, error) {
	var conditions []string
	var args []interface{}
	argCount := 0

	// Build WHERE clause
	conditions = append(conditions, "user_id = $1")
	args = append(args, filter.UserID)
	argCount = 1

	if filter.AccountID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("account_id = $%d", argCount))
		args = append(args, *filter.AccountID)
	}

	if filter.CategoryID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argCount))
		args = append(args, *filter.CategoryID)
	}

	if filter.Type != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("type = $%d", argCount))
		args = append(args, *filter.Type)
	}

	if filter.StartDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("transaction_date >= $%d", argCount))
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("transaction_date <= $%d", argCount))
		args = append(args, *filter.EndDate)
	}

	// เพิ่มการค้นหาจาก note หรือชื่อหมวดหมู่
	if filter.SearchQuery != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf(`(
			note ILIKE $%d OR 
			EXISTS (
				SELECT 1 FROM categories c 
				WHERE c.id = transactions.category_id 
				AND c.name ILIKE $%d
			)
		)`, argCount, argCount))
		searchPattern := "%" + *filter.SearchQuery + "%"
		args = append(args, searchPattern)
	}

	// เพิ่มการกรองตามจำนวนเงิน
	if filter.MinAmount != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("amount >= $%d", argCount))
		args = append(args, *filter.MinAmount)
	}

	if filter.MaxAmount != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("amount <= $%d", argCount))
		args = append(args, *filter.MaxAmount)
	}

	query := `
		SELECT t.id, t.user_id, t.category_id, t.account_id, t.amount, t.type, t.note, t.transaction_date, t.created_at, t.updated_at
		FROM transactions t
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY t.transaction_date DESC, t.created_at DESC
	`

	if filter.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entity.Transaction
	for rows.Next() {
		transaction := &entity.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.Amount,
			&transaction.Type,
			&transaction.Note,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	query := `
		UPDATE transactions 
		SET category_id = $2, account_id = $3, amount = $4, type = $5, note = $6, transaction_date = $7, updated_at = $8
		WHERE id = $1
	`

	transaction.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		transaction.ID,
		transaction.CategoryID,
		transaction.AccountID,
		transaction.Amount,
		transaction.Type,
		transaction.Note,
		transaction.TransactionDate,
		transaction.UpdatedAt,
	)

	return err
}

func (r *transactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM transactions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *transactionRepository) GetMonthlySpending(ctx context.Context, userID uuid.UUID, year int, month int) (map[uuid.UUID]float64, error) {
	query := `
		SELECT category_id, SUM(amount::numeric) as total
		FROM transactions 
		WHERE user_id = $1 
		AND type = 'expense'
		AND EXTRACT(YEAR FROM transaction_date) = $2
		AND EXTRACT(MONTH FROM transaction_date) = $3
		GROUP BY category_id
	`

	rows, err := r.db.QueryContext(ctx, query, userID, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]float64)
	for rows.Next() {
		var categoryID uuid.UUID
		var total decimal.Decimal

		err := rows.Scan(&categoryID, &total)
		if err != nil {
			return nil, err
		}

		floatTotal, _ := total.Float64()
		result[categoryID] = floatTotal
	}

	return result, nil
}
