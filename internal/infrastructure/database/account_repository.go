package database

import (
	"context"
	"database/sql"
	"time"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"

	"github.com/google/uuid"
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) repository.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *entity.Account) error {
	query := `
		INSERT INTO accounts (id, user_id, name, type, initial_balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.UserID,
		account.Name,
		account.Type,
		account.InitialBalance,
		account.CreatedAt,
		account.UpdatedAt,
	)

	return err
}

func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	query := `
		SELECT id, user_id, name, type, initial_balance, created_at, updated_at
		FROM accounts WHERE id = $1
	`

	account := &entity.Account{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.Type,
		&account.InitialBalance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error) {
	query := `
		SELECT id, user_id, name, type, initial_balance, created_at, updated_at
		FROM accounts 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*entity.Account
	for rows.Next() {
		account := &entity.Account{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Name,
			&account.Type,
			&account.InitialBalance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *accountRepository) Update(ctx context.Context, account *entity.Account) error {
	query := `
		UPDATE accounts 
		SET name = $2, type = $3, initial_balance = $4, updated_at = $5
		WHERE id = $1
	`

	account.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.Name,
		account.Type,
		account.InitialBalance,
		account.UpdatedAt,
	)

	return err
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
