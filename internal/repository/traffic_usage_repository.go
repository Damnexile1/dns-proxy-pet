package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/damnexile/dns-proxy-pet/internal/database"
	"github.com/damnexile/dns-proxy-pet/internal/models"
	"github.com/jackc/pgx/v5"
)

// TrafficUsageRepository handles database operations for traffic usage
type TrafficUsageRepository struct {
	db database.Database
}

// NewTrafficUsageRepository creates a new traffic usage repository
func NewTrafficUsageRepository(db database.Database) *TrafficUsageRepository {
	return &TrafficUsageRepository{db: db}
}

// Create creates a new traffic usage record
func (r *TrafficUsageRepository) Create(ctx context.Context, traffic *models.TrafficUsage) error {
	query := `
		INSERT INTO traffic_usage (user_id, bytes_used, date, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, traffic.UserID, traffic.BytesUsed, traffic.Date).
		Scan(&traffic.ID, &traffic.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create traffic usage: %w", err)
	}

	return nil
}

// GetByID retrieves traffic usage by ID
func (r *TrafficUsageRepository) GetByID(ctx context.Context, id int64) (*models.TrafficUsage, error) {
	query := `
		SELECT id, user_id, bytes_used, date, created_at
		FROM traffic_usage
		WHERE id = $1
	`

	var traffic models.TrafficUsage
	err := r.db.QueryRow(ctx, query, id).
		Scan(&traffic.ID, &traffic.UserID, &traffic.BytesUsed, &traffic.Date, &traffic.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("traffic usage not found")
		}
		return nil, fmt.Errorf("failed to get traffic usage: %w", err)
	}

	return &traffic, nil
}

// GetByUserAndDate retrieves traffic usage for a user on a specific date
func (r *TrafficUsageRepository) GetByUserAndDate(ctx context.Context, userID int64, date time.Time) (*models.TrafficUsage, error) {
	query := `
		SELECT id, user_id, bytes_used, date, created_at
		FROM traffic_usage
		WHERE user_id = $1 AND date = $2
	`

	var traffic models.TrafficUsage
	err := r.db.QueryRow(ctx, query, userID, date).
		Scan(&traffic.ID, &traffic.UserID, &traffic.BytesUsed, &traffic.Date, &traffic.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("traffic usage not found")
		}
		return nil, fmt.Errorf("failed to get traffic usage: %w", err)
	}

	return &traffic, nil
}

// IncrementUsage increments traffic usage for a user on a specific date
func (r *TrafficUsageRepository) IncrementUsage(ctx context.Context, userID int64, date time.Time, bytes int64) error {
	query := `
		INSERT INTO traffic_usage (user_id, bytes_used, date, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, date)
		DO UPDATE SET bytes_used = traffic_usage.bytes_used + $2
	`

	_, err := r.db.Exec(ctx, query, userID, bytes, date)
	if err != nil {
		return fmt.Errorf("failed to increment traffic usage: %w", err)
	}

	return nil
}

// GetByUserIDAndDateRange retrieves traffic usage for a user within a date range
func (r *TrafficUsageRepository) GetByUserIDAndDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*models.TrafficUsage, error) {
	query := `
		SELECT id, user_id, bytes_used, date, created_at
		FROM traffic_usage
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get traffic usage: %w", err)
	}
	defer rows.Close()

	var usages []*models.TrafficUsage
	for rows.Next() {
		var traffic models.TrafficUsage
		err := rows.Scan(&traffic.ID, &traffic.UserID, &traffic.BytesUsed, &traffic.Date, &traffic.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan traffic usage: %w", err)
		}
		usages = append(usages, &traffic)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating traffic usage: %w", err)
	}

	return usages, nil
}

// GetTotalByUserID retrieves total traffic usage for a user
func (r *TrafficUsageRepository) GetTotalByUserID(ctx context.Context, userID int64) (int64, error) {
	query := `
		SELECT COALESCE(SUM(bytes_used), 0)
		FROM traffic_usage
		WHERE user_id = $1
	`

	var total int64
	err := r.db.QueryRow(ctx, query, userID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total traffic usage: %w", err)
	}

	return total, nil
}

// GetTotalByUserIDAndDateRange retrieves total traffic usage for a user within a date range
func (r *TrafficUsageRepository) GetTotalByUserIDAndDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) (int64, error) {
	query := `
		SELECT COALESCE(SUM(bytes_used), 0)
		FROM traffic_usage
		WHERE user_id = $1 AND date >= $2 AND date <= $3
	`

	var total int64
	err := r.db.QueryRow(ctx, query, userID, startDate, endDate).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total traffic usage: %w", err)
	}

	return total, nil
}

// Delete deletes traffic usage record
func (r *TrafficUsageRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM traffic_usage WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete traffic usage: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("traffic usage not found")
	}

	return nil
}
