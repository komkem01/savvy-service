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

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	query := `
		INSERT INTO categories (id, user_id, name, type, icon_name, color_hex, is_archived, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		category.ID,
		category.UserID,
		category.Name,
		category.Type,
		category.IconName,
		category.ColorHex,
		category.IsArchived,
		category.CreatedAt,
		category.UpdatedAt,
	)

	return err
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	query := `
		SELECT id, user_id, name, type, icon_name, color_hex, is_archived, created_at, updated_at
		FROM categories WHERE id = $1 AND is_archived = false
	`

	category := &entity.Category{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.Type,
		&category.IconName,
		&category.ColorHex,
		&category.IsArchived,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error) {
	query := `
		SELECT id, user_id, name, type, icon_name, color_hex, is_archived, created_at, updated_at
		FROM categories 
		WHERE user_id = $1 AND is_archived = false
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		category := &entity.Category{}
		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.Type,
			&category.IconName,
			&category.ColorHex,
			&category.IsArchived,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *categoryRepository) GetSystemCategories(ctx context.Context) ([]*entity.Category, error) {
	query := `
		SELECT id, user_id, name, type, icon_name, color_hex, is_archived, created_at, updated_at
		FROM categories 
		WHERE user_id IS NULL AND is_archived = false
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		category := &entity.Category{}
		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.Type,
			&category.IconName,
			&category.ColorHex,
			&category.IsArchived,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	query := `
		UPDATE categories 
		SET name = $2, type = $3, icon_name = $4, color_hex = $5, updated_at = $6
		WHERE id = $1
	`

	category.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		category.ID,
		category.Name,
		category.Type,
		category.IconName,
		category.ColorHex,
		category.UpdatedAt,
	)

	return err
}

func (r *categoryRepository) GetByFilter(ctx context.Context, filter repository.CategoryFilter) ([]*entity.Category, error) {
	var conditions []string
	var args []interface{}
	argCount := 0

	// Build WHERE clause
	if filter.UserID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, *filter.UserID)
	}

	if filter.IsSystem != nil {
		if *filter.IsSystem {
			conditions = append(conditions, "user_id IS NULL")
		} else {
			conditions = append(conditions, "user_id IS NOT NULL")
		}
	}

	if filter.IsArchived != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("is_archived = $%d", argCount))
		args = append(args, *filter.IsArchived)
	}

	if filter.Type != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("type = $%d", argCount))
		args = append(args, *filter.Type)
	}

	if filter.SearchName != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argCount))
		searchPattern := "%" + *filter.SearchName + "%"
		args = append(args, searchPattern)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, name, type, icon_name, color_hex, is_archived, created_at, updated_at
		FROM categories 
		%s
		ORDER BY name ASC
	`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		category := &entity.Category{}
		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.Type,
			&category.IconName,
			&category.ColorHex,
			&category.IsArchived,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *categoryRepository) Unarchive(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE categories SET is_archived = false, updated_at = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, time.Now())
	return err
}

func (r *categoryRepository) GetCategoryUsageStats(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]int64, error) {
	query := `
		SELECT c.id, COUNT(t.id) as usage_count
		FROM categories c
		LEFT JOIN transactions t ON c.id = t.category_id
		WHERE (c.user_id = $1 OR c.user_id IS NULL)
		GROUP BY c.id
		ORDER BY usage_count DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]int64)
	for rows.Next() {
		var categoryID uuid.UUID
		var usageCount int64

		err := rows.Scan(&categoryID, &usageCount)
		if err != nil {
			return nil, err
		}

		result[categoryID] = usageCount
	}

	return result, nil
}

func (r *categoryRepository) Archive(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE categories SET is_archived = true, updated_at = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, time.Now())
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
