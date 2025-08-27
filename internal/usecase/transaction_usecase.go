package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"
	"savvy-backend/pkg/utils"
)

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, userID, categoryID, accountID uuid.UUID,
		amount decimal.Decimal, transactionType entity.TransactionType,
		note *string, transactionDate string) (*entity.Transaction, error)
	GetTransactionsByFilter(ctx context.Context, filter repository.TransactionFilter) ([]*entity.Transaction, error)
	GetTransactionByID(ctx context.Context, userID, transactionID uuid.UUID) (*entity.Transaction, error)
	UpdateTransaction(ctx context.Context, userID uuid.UUID, transaction *entity.Transaction) error
	DeleteTransaction(ctx context.Context, userID, transactionID uuid.UUID) error
	GetMonthlyReport(ctx context.Context, userID uuid.UUID, year, month int) (map[string]interface{}, error)
}

type transactionUsecase struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	categoryRepo    repository.CategoryRepository
}

func NewTransactionUsecase(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
	}
}

func (t *transactionUsecase) CreateTransaction(ctx context.Context, userID, categoryID, accountID uuid.UUID,
	amount decimal.Decimal, transactionType entity.TransactionType,
	note *string, transactionDate string) (*entity.Transaction, error) {

	// Validate account belongs to user
	account, err := t.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, errors.New("account not found")
	}
	if account.UserID != userID {
		return nil, errors.New("account does not belong to user")
	}

	// Validate category
	category, err := t.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}
	if category.UserID != nil && *category.UserID != userID {
		return nil, errors.New("category does not belong to user")
	}

	// Parse transaction date
	parsedDate, err := utils.ParseDate(transactionDate)
	if err != nil {
		return nil, errors.New("invalid transaction date format")
	}

	// Create transaction
	transaction := entity.NewTransaction(userID, categoryID, accountID, amount, transactionType, note, parsedDate)

	err = t.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *transactionUsecase) GetTransactionsByFilter(ctx context.Context, filter repository.TransactionFilter) ([]*entity.Transaction, error) {
	return t.transactionRepo.GetByFilter(ctx, filter)
}

func (t *transactionUsecase) GetTransactionByID(ctx context.Context, userID, transactionID uuid.UUID) (*entity.Transaction, error) {
	transaction, err := t.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	if transaction.UserID != userID {
		return nil, errors.New("transaction does not belong to user")
	}

	return transaction, nil
}

func (t *transactionUsecase) UpdateTransaction(ctx context.Context, userID uuid.UUID, transaction *entity.Transaction) error {
	// Verify ownership
	if transaction.UserID != userID {
		return errors.New("transaction does not belong to user")
	}

	return t.transactionRepo.Update(ctx, transaction)
}

func (t *transactionUsecase) DeleteTransaction(ctx context.Context, userID, transactionID uuid.UUID) error {
	transaction, err := t.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return err
	}

	if transaction.UserID != userID {
		return errors.New("transaction does not belong to user")
	}

	return t.transactionRepo.Delete(ctx, transactionID)
}

func (t *transactionUsecase) GetMonthlyReport(ctx context.Context, userID uuid.UUID, year, month int) (map[string]interface{}, error) {
	spendingByCategory, err := t.transactionRepo.GetMonthlySpending(ctx, userID, year, month)
	if err != nil {
		return nil, err
	}

	// Additional logic for monthly report can be added here

	return map[string]interface{}{
		"spending_by_category": spendingByCategory,
		"year":                 year,
		"month":                month,
	}, nil
}
