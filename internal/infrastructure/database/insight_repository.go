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

type insightRepository struct {
	db *sql.DB
}

func NewInsightRepository(db *sql.DB) repository.InsightRepository {
	return &insightRepository{db: db}
}

func (r *insightRepository) Create(ctx context.Context, insight *entity.Insight) error {
	query := `
		INSERT INTO insights (
			id, user_id, type, priority, title, content, action_text, is_read,
			related_entity_id, related_entity_type, related_data, valid_until,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.ExecContext(ctx, query,
		insight.ID,
		insight.UserID,
		insight.Type,
		insight.Priority,
		insight.Title,
		insight.Content,
		insight.ActionText,
		insight.IsRead,
		insight.RelatedEntityID,
		insight.RelatedEntityType,
		insight.RelatedData,
		insight.ValidUntil,
		insight.CreatedAt,
		insight.UpdatedAt,
	)

	return err
}

func (r *insightRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Insight, error) {
	query := `
		SELECT id, user_id, type, priority, title, content, action_text, is_read,
			   related_entity_id, related_entity_type, related_data, valid_until,
			   created_at, updated_at
		FROM insights WHERE id = $1
	`

	insight := &entity.Insight{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&insight.ID,
		&insight.UserID,
		&insight.Type,
		&insight.Priority,
		&insight.Title,
		&insight.Content,
		&insight.ActionText,
		&insight.IsRead,
		&insight.RelatedEntityID,
		&insight.RelatedEntityType,
		&insight.RelatedData,
		&insight.ValidUntil,
		&insight.CreatedAt,
		&insight.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return insight, nil
}

func (r *insightRepository) GetByFilter(ctx context.Context, filter repository.InsightFilter) ([]*entity.Insight, error) {
	var conditions []string
	var args []interface{}
	argCount := 0

	// Build WHERE clause
	conditions = append(conditions, "user_id = $1")
	args = append(args, filter.UserID)
	argCount = 1

	if filter.Type != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("type = $%d", argCount))
		args = append(args, *filter.Type)
	}

	if filter.Priority != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argCount))
		args = append(args, *filter.Priority)
	}

	if filter.IsRead != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("is_read = $%d", argCount))
		args = append(args, *filter.IsRead)
	}

	if filter.ValidOnly {
		conditions = append(conditions, "(valid_until IS NULL OR valid_until > NOW())")
	}

	query := `
		SELECT id, user_id, type, priority, title, content, action_text, is_read,
			   related_entity_id, related_entity_type, related_data, valid_until,
			   created_at, updated_at
		FROM insights 
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY priority DESC, created_at DESC
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

	var insights []*entity.Insight
	for rows.Next() {
		insight := &entity.Insight{}
		err := rows.Scan(
			&insight.ID,
			&insight.UserID,
			&insight.Type,
			&insight.Priority,
			&insight.Title,
			&insight.Content,
			&insight.ActionText,
			&insight.IsRead,
			&insight.RelatedEntityID,
			&insight.RelatedEntityType,
			&insight.RelatedData,
			&insight.ValidUntil,
			&insight.CreatedAt,
			&insight.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		insights = append(insights, insight)
	}

	return insights, nil
}

func (r *insightRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.Insight, error) {
	filter := repository.InsightFilter{
		UserID:    userID,
		ValidOnly: true,
		Limit:     limit,
	}
	return r.GetByFilter(ctx, filter)
}

func (r *insightRepository) GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Insight, error) {
	isRead := false
	filter := repository.InsightFilter{
		UserID:    userID,
		IsRead:    &isRead,
		ValidOnly: true,
		Limit:     50,
	}
	return r.GetByFilter(ctx, filter)
}

func (r *insightRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE insights SET is_read = true, updated_at = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, time.Now())
	return err
}

func (r *insightRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE insights SET is_read = true, updated_at = $2 WHERE user_id = $1 AND is_read = false`
	_, err := r.db.ExecContext(ctx, query, userID, time.Now())
	return err
}

func (r *insightRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM insights WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *insightRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	query := `DELETE FROM insights WHERE valid_until IS NOT NULL AND valid_until < $1`
	_, err := r.db.ExecContext(ctx, query, before)
	return err
}

func (r *insightRepository) GetSpendingAnomalies(ctx context.Context, userID uuid.UUID, months int) ([]*entity.SpendingAnomaly, error) {
	query := `
		WITH monthly_spending AS (
			SELECT 
				t.category_id,
				c.name as category_name,
				DATE_TRUNC('month', t.transaction_date) as month,
				SUM(t.amount) as total_amount
			FROM transactions t
			INNER JOIN categories c ON t.category_id = c.id
			WHERE t.user_id = $1 
			  AND t.type = 'expense'
			  AND t.transaction_date >= NOW() - INTERVAL '%d months'
			GROUP BY t.category_id, c.name, DATE_TRUNC('month', t.transaction_date)
		),
		averages AS (
			SELECT 
				category_id,
				category_name,
				AVG(total_amount) as avg_amount,
				STDDEV(total_amount) as stddev_amount
			FROM monthly_spending
			GROUP BY category_id, category_name
			HAVING COUNT(*) >= 3 -- Need at least 3 months of data
		),
		current_month AS (
			SELECT 
				t.category_id,
				c.name as category_name,
				SUM(t.amount) as current_amount
			FROM transactions t
			INNER JOIN categories c ON t.category_id = c.id
			WHERE t.user_id = $1 
			  AND t.type = 'expense'
			  AND DATE_TRUNC('month', t.transaction_date) = DATE_TRUNC('month', NOW())
			GROUP BY t.category_id, c.name
		)
		SELECT 
			cm.category_id,
			cm.category_name,
			cm.current_amount,
			a.avg_amount,
			ROUND(((cm.current_amount - a.avg_amount) / a.avg_amount * 100)::numeric, 2) as percentage_increase
		FROM current_month cm
		INNER JOIN averages a ON cm.category_id = a.category_id
		WHERE cm.current_amount > a.avg_amount + (a.stddev_amount * 1.5) -- 1.5 standard deviations
		  AND ((cm.current_amount - a.avg_amount) / a.avg_amount) > 0.3 -- At least 30% increase
		ORDER BY percentage_increase DESC
	`

	formattedQuery := fmt.Sprintf(query, months)
	rows, err := r.db.QueryContext(ctx, formattedQuery, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var anomalies []*entity.SpendingAnomaly
	for rows.Next() {
		anomaly := &entity.SpendingAnomaly{}
		err := rows.Scan(
			&anomaly.CategoryID,
			&anomaly.CategoryName,
			&anomaly.CurrentAmount,
			&anomaly.AverageAmount,
			&anomaly.PercentageIncrease,
		)
		if err != nil {
			return nil, err
		}

		// Determine severity
		if anomaly.PercentageIncrease >= 100 {
			anomaly.Severity = "high"
		} else if anomaly.PercentageIncrease >= 50 {
			anomaly.Severity = "medium"
		} else {
			anomaly.Severity = "low"
		}

		anomaly.Period = time.Now().Format("2006-01")
		anomalies = append(anomalies, anomaly)
	}

	return anomalies, nil
}

func (r *insightRepository) GetSpendingPatterns(ctx context.Context, userID uuid.UUID, days int) ([]*entity.SpendingPattern, error) {
	query := `
		SELECT 
			t.category_id,
			c.name as category_name,
			TO_CHAR(t.transaction_date, 'Day') as day_of_week,
			CASE 
				WHEN EXTRACT(HOUR FROM t.created_at) BETWEEN 6 AND 11 THEN 'Morning'
				WHEN EXTRACT(HOUR FROM t.created_at) BETWEEN 12 AND 17 THEN 'Afternoon'
				WHEN EXTRACT(HOUR FROM t.created_at) BETWEEN 18 AND 22 THEN 'Evening'
				ELSE 'Night'
			END as time_of_day,
			COUNT(*) as frequency_count,
			AVG(t.amount) as avg_amount
		FROM transactions t
		INNER JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = $1 
		  AND t.type = 'expense'
		  AND t.transaction_date >= NOW() - INTERVAL '%d days'
		GROUP BY t.category_id, c.name, TO_CHAR(t.transaction_date, 'Day'), 
				 CASE 
					WHEN EXTRACT(HOUR FROM t.created_at) BETWEEN 6 AND 11 THEN 'Morning'
					WHEN EXTRACT(HOUR FROM t.created_at) BETWEEN 12 AND 17 THEN 'Afternoon'
					WHEN EXTRACT(HOUR FROM t.created_at) BETWEEN 18 AND 22 THEN 'Evening'
					ELSE 'Night'
				 END
		HAVING COUNT(*) >= 3 -- At least 3 transactions in this pattern
		ORDER BY frequency_count DESC, avg_amount DESC
		LIMIT 20
	`

	formattedQuery := fmt.Sprintf(query, days)
	rows, err := r.db.QueryContext(ctx, formattedQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patterns []*entity.SpendingPattern
	for rows.Next() {
		pattern := &entity.SpendingPattern{}
		err := rows.Scan(
			&pattern.CategoryID,
			&pattern.CategoryName,
			&pattern.DayOfWeek,
			&pattern.TimeOfDay,
			&pattern.FrequencyCount,
			&pattern.AverageAmount,
		)
		if err != nil {
			return nil, err
		}

		pattern.DayOfWeek = strings.TrimSpace(pattern.DayOfWeek)
		pattern.Period = fmt.Sprintf("Last %d days", days)
		patterns = append(patterns, pattern)
	}

	return patterns, nil
}
