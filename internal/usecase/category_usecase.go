package usecase

import (
	"context"
	"errors"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"

	"github.com/google/uuid"
)

type CategoryUsecase interface {
	CreateCategory(ctx context.Context, userID uuid.UUID, name string, categoryType entity.CategoryType, iconName, colorHex *string) (*entity.Category, error)
	GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error)
	GetSystemCategories(ctx context.Context) ([]*entity.Category, error)
	GetAllAvailableCategories(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error)
	GetUserCategories(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error)
	GetCategoryUsageStats(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]int64, error)
	UpdateCategory(ctx context.Context, userID uuid.UUID, category *entity.Category) error
	ArchiveCategory(ctx context.Context, userID, categoryID uuid.UUID) error
	InitializeDefaultCategories(ctx context.Context) error
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (c *categoryUsecase) CreateCategory(ctx context.Context, userID uuid.UUID, name string, categoryType entity.CategoryType, iconName, colorHex *string) (*entity.Category, error) {
	category := entity.NewCategory(&userID, name, categoryType)
	category.IconName = iconName
	category.ColorHex = colorHex

	err := c.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c *categoryUsecase) GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error) {
	return c.categoryRepo.GetByUserID(ctx, userID)
}

func (c *categoryUsecase) GetSystemCategories(ctx context.Context) ([]*entity.Category, error) {
	return c.categoryRepo.GetSystemCategories(ctx)
}

func (c *categoryUsecase) GetAllAvailableCategories(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error) {
	// Get both system categories and user categories
	systemCategories, err := c.categoryRepo.GetSystemCategories(ctx)
	if err != nil {
		return nil, err
	}

	userCategories, err := c.categoryRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Combine both lists
	allCategories := append(systemCategories, userCategories...)
	return allCategories, nil
}

func (c *categoryUsecase) UpdateCategory(ctx context.Context, userID uuid.UUID, category *entity.Category) error {
	// Verify ownership (only user categories can be updated)
	if category.UserID == nil || *category.UserID != userID {
		return errors.New("category does not belong to user or is a system category")
	}

	return c.categoryRepo.Update(ctx, category)
}

func (c *categoryUsecase) ArchiveCategory(ctx context.Context, userID, categoryID uuid.UUID) error {
	category, err := c.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return err
	}

	if category.UserID == nil || *category.UserID != userID {
		return errors.New("category does not belong to user or is a system category")
	}

	return c.categoryRepo.Archive(ctx, categoryID)
}

func (c *categoryUsecase) InitializeDefaultCategories(ctx context.Context) error {
	// Default expense categories
	expenseCategories := []struct {
		name     string
		iconName string
		colorHex string
	}{
		{"อาหาร", "🍽️", "#FF6B6B"},
		{"เดินทาง", "🚗", "#4ECDC4"},
		{"ช้อปปิ้ง", "🛍️", "#45B7D1"},
		{"บันเทิง", "🎬", "#96CEB4"},
		{"สุขภาพ", "🏥", "#FECA57"},
		{"การศึกษา", "📚", "#FF9FF3"},
		{"บิล/ค่าใช้จ่าย", "💡", "#54A0FF"},
		{"อื่นๆ", "📝", "#5F27CD"},
	}

	// Default income categories
	incomeCategories := []struct {
		name     string
		iconName string
		colorHex string
	}{
		{"เงินเดือน", "💰", "#00D2D3"},
		{"โบนัส", "🎁", "#FF9F43"},
		{"ลงทุน", "📈", "#5f27cd"},
		{"ธุรกิจ", "🏢", "#00d2d3"},
		{"อื่นๆ", "💵", "#2ed573"},
	}

	// Create expense categories
	for _, cat := range expenseCategories {
		iconName := cat.iconName
		colorHex := cat.colorHex
		category := entity.NewCategory(nil, cat.name, entity.CategoryTypeExpense)
		category.IconName = &iconName
		category.ColorHex = &colorHex

		err := c.categoryRepo.Create(ctx, category)
		if err != nil {
			return err
		}
	}

	// Create income categories
	for _, cat := range incomeCategories {
		iconName := cat.iconName
		colorHex := cat.colorHex
		category := entity.NewCategory(nil, cat.name, entity.CategoryTypeIncome)
		category.IconName = &iconName
		category.ColorHex = &colorHex

		err := c.categoryRepo.Create(ctx, category)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *categoryUsecase) GetUserCategories(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error) {
	return c.GetAllAvailableCategories(ctx, userID)
}

func (c *categoryUsecase) GetCategoryUsageStats(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]int64, error) {
	return c.categoryRepo.GetCategoryUsageStats(ctx, userID)
}
