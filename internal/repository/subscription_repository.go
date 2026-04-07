package repository

import (
	"context"
	"fmt"

	"github.com/damnexile/dns-proxy-pet/internal/database"
	"github.com/damnexile/dns-proxy-pet/internal/models"
	"github.com/jackc/pgx/v5"
)

// SubscriptionRepository handles database operations for subscriptions
type SubscriptionRepository struct {
	db database.Database
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db database.Database) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, plan_type, status, started_at, expires_at, auto_renew, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		sub.UserID, sub.PlanType, sub.Status, sub.StartedAt, sub.ExpiresAt, sub.AutoRenew).
		Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

// GetByID retrieves a subscription by ID
func (r *SubscriptionRepository) GetByID(ctx context.Context, id int64) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_type, status, started_at, expires_at, auto_renew, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var sub models.Subscription
	err := r.db.QueryRow(ctx, query, id).
		Scan(&sub.ID, &sub.UserID, &sub.PlanType, &sub.Status, &sub.StartedAt, &sub.ExpiresAt, &sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &sub, nil
}

// GetByUserID retrieves active subscription for a user
func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID int64) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_type, status, started_at, expires_at, auto_renew, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var sub models.Subscription
	err := r.db.QueryRow(ctx, query, userID, models.SubscriptionStatusActive).
		Scan(&sub.ID, &sub.UserID, &sub.PlanType, &sub.Status, &sub.StartedAt, &sub.ExpiresAt, &sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("active subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &sub, nil
}

// Update updates a subscription
func (r *SubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	query := `
		UPDATE subscriptions
		SET plan_type = $1, status = $2, expires_at = $3, auto_renew = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query, sub.PlanType, sub.Status, sub.ExpiresAt, sub.AutoRenew, sub.ID).
		Scan(&sub.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// UpdateStatus updates subscription status
func (r *SubscriptionRepository) UpdateStatus(ctx context.Context, id int64, status models.SubscriptionStatus) error {
	query := `
		UPDATE subscriptions
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update subscription status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// Delete deletes a subscription
func (r *SubscriptionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// ListByUserID retrieves all subscriptions for a user
func (r *SubscriptionRepository) ListByUserID(ctx context.Context, userID int64) ([]*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_type, status, started_at, expires_at, auto_renew, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(&sub.ID, &sub.UserID, &sub.PlanType, &sub.Status, &sub.StartedAt, &sub.ExpiresAt, &sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subscriptions, nil
}

// GetExpiring retrieves subscriptions expiring within the given days
func (r *SubscriptionRepository) GetExpiring(ctx context.Context, days int) ([]*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_type, status, started_at, expires_at, auto_renew, created_at, updated_at
		FROM subscriptions
		WHERE status = $1 AND expires_at <= NOW() + INTERVAL '1 day' * $2
		ORDER BY expires_at ASC
	`

	rows, err := r.db.Query(ctx, query, models.SubscriptionStatusActive, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(&sub.ID, &sub.UserID, &sub.PlanType, &sub.Status, &sub.StartedAt, &sub.ExpiresAt, &sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subscriptions, nil
}
