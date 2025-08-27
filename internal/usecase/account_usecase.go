package usecase

import (
	"context"
	"errors"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountUsecase interface {
	CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountType entity.AccountType) (*entity.Account, error)
	GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error)
	GetAccountByID(ctx context.Context, userID, accountID uuid.UUID) (*entity.Account, error)
	UpdateAccount(ctx context.Context, userID uuid.UUID, account *entity.Account) error
	DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error
}

type accountUsecase struct {
	accountRepo repository.AccountRepository
}

func NewAccountUsecase(accountRepo repository.AccountRepository) AccountUsecase {
	return &accountUsecase{
		accountRepo: accountRepo,
	}
}

func (a *accountUsecase) CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountType entity.AccountType) (*entity.Account, error) {
	// ใส่ initial balance เป็น 0
	account := entity.NewAccount(userID, name, accountType, decimal.Zero)

	err := a.accountRepo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountUsecase) GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error) {
	return a.accountRepo.GetByUserID(ctx, userID)
}

func (a *accountUsecase) GetAccountByID(ctx context.Context, userID, accountID uuid.UUID) (*entity.Account, error) {
	account, err := a.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if account.UserID != userID {
		return nil, errors.New("account does not belong to user")
	}

	return account, nil
}

func (a *accountUsecase) UpdateAccount(ctx context.Context, userID uuid.UUID, account *entity.Account) error {
	if account.UserID != userID {
		return errors.New("account does not belong to user")
	}

	return a.accountRepo.Update(ctx, account)
}

func (a *accountUsecase) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	account, err := a.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	if account.UserID != userID {
		return errors.New("account does not belong to user")
	}

	return a.accountRepo.Delete(ctx, accountID)
}
