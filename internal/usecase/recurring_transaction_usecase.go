package usecase

import (
	"context"
	"fmt"
	"time"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"

	"github.com/google/uuid"
)

type RecurringTransactionUsecase interface {
	CreateRecurringTransaction(ctx context.Context, userID uuid.UUID, recurring *entity.RecurringTransaction) (*entity.RecurringTransaction, error)
	GetUserRecurringTransactions(ctx context.Context, userID uuid.UUID) ([]*entity.RecurringTransaction, error)
	GetRecurringTransactionByID(ctx context.Context, userID, recurringID uuid.UUID) (*entity.RecurringTransaction, error)
	UpdateRecurringTransaction(ctx context.Context, userID uuid.UUID, recurring *entity.RecurringTransaction) error
	DeleteRecurringTransaction(ctx context.Context, userID, recurringID uuid.UUID) error
	GetDueTransactions(ctx context.Context, userID uuid.UUID) ([]*entity.RecurringTransaction, error)
	ExecuteRecurringTransaction(ctx context.Context, userID, recurringID uuid.UUID) (*entity.Transaction, error)
	ProcessAllDueTransactions(ctx context.Context) error
}

type recurringTransactionUsecase struct {
	recurringRepo   repository.RecurringTransactionRepository
	transactionRepo repository.TransactionRepository
	categoryRepo    repository.CategoryRepository
	accountRepo     repository.AccountRepository
}

func NewRecurringTransactionUsecase(
	recurringRepo repository.RecurringTransactionRepository,
	transactionRepo repository.TransactionRepository,
	categoryRepo repository.CategoryRepository,
	accountRepo repository.AccountRepository,
) RecurringTransactionUsecase {
	return &recurringTransactionUsecase{
		recurringRepo:   recurringRepo,
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
		accountRepo:     accountRepo,
	}
}

func (r *recurringTransactionUsecase) CreateRecurringTransaction(ctx context.Context, userID uuid.UUID, recurring *entity.RecurringTransaction) (*entity.RecurringTransaction, error) {
	// Validate ownership
	recurring.UserID = userID

	// Validate category belongs to user
	category, err := r.categoryRepo.GetByID(ctx, recurring.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	if category.UserID != nil && *category.UserID != userID {
		return nil, fmt.Errorf("category does not belong to user")
	}

	// Validate account belongs to user
	account, err := r.accountRepo.GetByID(ctx, recurring.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	if account.UserID != userID {
		return nil, fmt.Errorf("account does not belong to user")
	}

	err = r.recurringRepo.Create(ctx, recurring)
	if err != nil {
		return nil, fmt.Errorf("failed to create recurring transaction: %w", err)
	}

	return recurring, nil
}

func (r *recurringTransactionUsecase) GetUserRecurringTransactions(ctx context.Context, userID uuid.UUID) ([]*entity.RecurringTransaction, error) {
	filter := repository.RecurringTransactionFilter{
		UserID: userID,
	}

	return r.recurringRepo.GetByFilter(ctx, filter)
}

func (r *recurringTransactionUsecase) GetRecurringTransactionByID(ctx context.Context, userID, recurringID uuid.UUID) (*entity.RecurringTransaction, error) {
	recurring, err := r.recurringRepo.GetByID(ctx, recurringID)
	if err != nil {
		return nil, err
	}

	if recurring.UserID != userID {
		return nil, fmt.Errorf("recurring transaction does not belong to user")
	}

	return recurring, nil
}

func (r *recurringTransactionUsecase) UpdateRecurringTransaction(ctx context.Context, userID uuid.UUID, recurring *entity.RecurringTransaction) error {
	existing, err := r.recurringRepo.GetByID(ctx, recurring.ID)
	if err != nil {
		return err
	}

	if existing.UserID != userID {
		return fmt.Errorf("recurring transaction does not belong to user")
	}

	recurring.UserID = userID // Ensure user ID is preserved
	return r.recurringRepo.Update(ctx, recurring)
}

func (r *recurringTransactionUsecase) DeleteRecurringTransaction(ctx context.Context, userID, recurringID uuid.UUID) error {
	recurring, err := r.recurringRepo.GetByID(ctx, recurringID)
	if err != nil {
		return err
	}

	if recurring.UserID != userID {
		return fmt.Errorf("recurring transaction does not belong to user")
	}

	return r.recurringRepo.Delete(ctx, recurringID)
}

func (r *recurringTransactionUsecase) GetDueTransactions(ctx context.Context, userID uuid.UUID) ([]*entity.RecurringTransaction, error) {
	now := time.Now()
	filter := repository.RecurringTransactionFilter{
		UserID:   userID,
		DueDate:  &now,
		IsActive: &[]bool{true}[0],
	}

	return r.recurringRepo.GetByFilter(ctx, filter)
}

func (r *recurringTransactionUsecase) ExecuteRecurringTransaction(ctx context.Context, userID, recurringID uuid.UUID) (*entity.Transaction, error) {
	recurring, err := r.GetRecurringTransactionByID(ctx, userID, recurringID)
	if err != nil {
		return nil, err
	}

	if !recurring.IsActive {
		return nil, fmt.Errorf("recurring transaction is not active")
	}

	if recurring.NextExecutionDate.After(time.Now()) {
		return nil, fmt.Errorf("transaction is not due yet")
	}

	// Create the actual transaction
	transaction := &entity.Transaction{
		ID:              uuid.New(),
		UserID:          recurring.UserID,
		CategoryID:      recurring.CategoryID,
		AccountID:       recurring.AccountID,
		Amount:          recurring.Amount,
		Type:            recurring.Type,
		Note:            recurring.Note,
		TransactionDate: time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = r.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update recurring transaction
	executedAt := time.Now()
	err = r.recurringRepo.MarkAsExecuted(ctx, recurringID, executedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to mark as executed: %w", err)
	}

	// Calculate next execution date
	nextDate := recurring.CalculateNextExecutionDate()

	// Check if we should deactivate based on end date or remaining executions
	shouldDeactivate := false
	if recurring.EndDate != nil && nextDate.After(*recurring.EndDate) {
		shouldDeactivate = true
	}
	if recurring.RemainingExecutions != nil && *recurring.RemainingExecutions <= 1 {
		shouldDeactivate = true
	}

	if shouldDeactivate {
		recurring.IsActive = false
	} else {
		recurring.NextExecutionDate = nextDate
	}

	err = r.recurringRepo.Update(ctx, recurring)
	if err != nil {
		return nil, fmt.Errorf("failed to update recurring transaction: %w", err)
	}

	return transaction, nil
}

func (r *recurringTransactionUsecase) ProcessAllDueTransactions(ctx context.Context) error {
	// This method is typically called by a cron job or scheduled task
	now := time.Now()
	dueTransactions, err := r.recurringRepo.GetDueTransactions(ctx, now)
	if err != nil {
		return fmt.Errorf("failed to get due transactions: %w", err)
	}

	var errors []error
	for _, recurring := range dueTransactions {
		if recurring.AutoExecute {
			_, err := r.ExecuteRecurringTransaction(ctx, recurring.UserID, recurring.ID)
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to execute recurring transaction %s: %w", recurring.ID, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("some recurring transactions failed to execute: %v", errors)
	}

	return nil
}
