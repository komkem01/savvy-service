package usecase

import (
	"context"
	"time"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"
	"savvy-backend/pkg/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DashboardUsecase interface {
	GetMonthlySummary(ctx context.Context, userID uuid.UUID, year, month int) (*MonthlySummary, error)
	GetCurrentMonthlySummary(ctx context.Context, userID uuid.UUID) (*MonthlySummary, error)
	GetRecentTransactions(ctx context.Context, userID uuid.UUID, limit int) ([]*TransactionWithDetails, error)
	GetSpendingByCategory(ctx context.Context, userID uuid.UUID, year, month int) ([]*CategorySpending, error)
}

type MonthlySummary struct {
	Year         int             `json:"year"`
	Month        int             `json:"month"`
	TotalIncome  decimal.Decimal `json:"total_income"`
	TotalExpense decimal.Decimal `json:"total_expense"`
	Balance      decimal.Decimal `json:"balance"`
	StartDate    time.Time       `json:"start_date"`
	EndDate      time.Time       `json:"end_date"`
}

type TransactionWithDetails struct {
	Transaction  *entity.Transaction `json:"transaction"`
	CategoryName string              `json:"category_name"`
	AccountName  string              `json:"account_name"`
}

type CategorySpending struct {
	CategoryID   uuid.UUID       `json:"category_id"`
	CategoryName string          `json:"category_name"`
	Amount       decimal.Decimal `json:"amount"`
	IconName     *string         `json:"icon_name,omitempty"`
	ColorHex     *string         `json:"color_hex,omitempty"`
}

type dashboardUsecase struct {
	transactionRepo repository.TransactionRepository
	categoryRepo    repository.CategoryRepository
	accountRepo     repository.AccountRepository
}

func NewDashboardUsecase(
	transactionRepo repository.TransactionRepository,
	categoryRepo repository.CategoryRepository,
	accountRepo repository.AccountRepository,
) DashboardUsecase {
	return &dashboardUsecase{
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
		accountRepo:     accountRepo,
	}
}

func (d *dashboardUsecase) GetMonthlySummary(ctx context.Context, userID uuid.UUID, year, month int) (*MonthlySummary, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := utils.EndOfMonth(startDate)

	// Get all transactions for the month
	filter := repository.TransactionFilter{
		UserID:    userID,
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	transactions, err := d.transactionRepo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	var totalIncome, totalExpense decimal.Decimal

	for _, transaction := range transactions {
		if transaction.Type == entity.TransactionTypeIncome {
			totalIncome = totalIncome.Add(transaction.Amount)
		} else if transaction.Type == entity.TransactionTypeExpense {
			totalExpense = totalExpense.Add(transaction.Amount)
		}
	}

	balance := totalIncome.Sub(totalExpense)

	return &MonthlySummary{
		Year:         year,
		Month:        month,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      balance,
		StartDate:    startDate,
		EndDate:      endDate,
	}, nil
}

func (d *dashboardUsecase) GetCurrentMonthlySummary(ctx context.Context, userID uuid.UUID) (*MonthlySummary, error) {
	now := time.Now()
	return d.GetMonthlySummary(ctx, userID, now.Year(), int(now.Month()))
}

func (d *dashboardUsecase) GetRecentTransactions(ctx context.Context, userID uuid.UUID, limit int) ([]*TransactionWithDetails, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	filter := repository.TransactionFilter{
		UserID: userID,
		Limit:  limit,
	}

	transactions, err := d.transactionRepo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	var result []*TransactionWithDetails

	for _, transaction := range transactions {
		// Get category details
		category, err := d.categoryRepo.GetByID(ctx, transaction.CategoryID)
		if err != nil {
			continue // Skip if category not found
		}

		// Get account details
		account, err := d.accountRepo.GetByID(ctx, transaction.AccountID)
		if err != nil {
			continue // Skip if account not found
		}

		result = append(result, &TransactionWithDetails{
			Transaction:  transaction,
			CategoryName: category.Name,
			AccountName:  account.Name,
		})
	}

	return result, nil
}

func (d *dashboardUsecase) GetSpendingByCategory(ctx context.Context, userID uuid.UUID, year, month int) ([]*CategorySpending, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := utils.EndOfMonth(startDate)

	// Get expense transactions for the month
	expenseType := entity.TransactionTypeExpense
	filter := repository.TransactionFilter{
		UserID:    userID,
		Type:      &expenseType,
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	transactions, err := d.transactionRepo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Group by category
	categoryTotals := make(map[uuid.UUID]decimal.Decimal)
	categoryNames := make(map[uuid.UUID]string)
	categoryIcons := make(map[uuid.UUID]*string)
	categoryColors := make(map[uuid.UUID]*string)

	for _, transaction := range transactions {
		categoryTotals[transaction.CategoryID] = categoryTotals[transaction.CategoryID].Add(transaction.Amount)

		// Get category details if not already cached
		if _, exists := categoryNames[transaction.CategoryID]; !exists {
			category, err := d.categoryRepo.GetByID(ctx, transaction.CategoryID)
			if err == nil {
				categoryNames[transaction.CategoryID] = category.Name
				categoryIcons[transaction.CategoryID] = category.IconName
				categoryColors[transaction.CategoryID] = category.ColorHex
			}
		}
	}

	var result []*CategorySpending
	for categoryID, amount := range categoryTotals {
		result = append(result, &CategorySpending{
			CategoryID:   categoryID,
			CategoryName: categoryNames[categoryID],
			Amount:       amount,
			IconName:     categoryIcons[categoryID],
			ColorHex:     categoryColors[categoryID],
		})
	}

	return result, nil
}
