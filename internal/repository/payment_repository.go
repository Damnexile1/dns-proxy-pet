package repository

import (
	"context"
	"fmt"

	"github.com/damnexile/dns-proxy-pet/internal/database"
	"github.com/damnexile/dns-proxy-pet/internal/models"
	"github.com/jackc/pgx/v5"
)

// PaymentRepository handles database operations for payments
type PaymentRepository struct {
	db database.Database
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db database.Database) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create creates a new payment
func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (user_id, subscription_id, amount, currency, payment_method, status, external_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		payment.UserID, payment.SubscriptionID, payment.Amount, payment.Currency,
		payment.PaymentMethod, payment.Status, payment.ExternalID).
		Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(ctx context.Context, id int64) (*models.Payment, error) {
	query := `
		SELECT id, user_id, subscription_id, amount, currency, payment_method, status, external_id, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	var payment models.Payment
	err := r.db.QueryRow(ctx, query, id).
		Scan(&payment.ID, &payment.UserID, &payment.SubscriptionID, &payment.Amount, &payment.Currency,
			&payment.PaymentMethod, &payment.Status, &payment.ExternalID, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}

// GetByExternalID retrieves a payment by external ID
func (r *PaymentRepository) GetByExternalID(ctx context.Context, externalID string) (*models.Payment, error) {
	query := `
		SELECT id, user_id, subscription_id, amount, currency, payment_method, status, external_id, created_at, updated_at
		FROM payments
		WHERE external_id = $1
	`

	var payment models.Payment
	err := r.db.QueryRow(ctx, query, externalID).
		Scan(&payment.ID, &payment.UserID, &payment.SubscriptionID, &payment.Amount, &payment.Currency,
			&payment.PaymentMethod, &payment.Status, &payment.ExternalID, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}

// UpdateStatus updates payment status
func (r *PaymentRepository) UpdateStatus(ctx context.Context, id int64, status models.PaymentStatus) error {
	query := `
		UPDATE payments
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// Update updates a payment
func (r *PaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	query := `
		UPDATE payments
		SET subscription_id = $1, amount = $2, currency = $3, payment_method = $4, status = $5, external_id = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		payment.SubscriptionID, payment.Amount, payment.Currency, payment.PaymentMethod,
		payment.Status, payment.ExternalID, payment.ID).
		Scan(&payment.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

// ListByUserID retrieves all payments for a user
func (r *PaymentRepository) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.Payment, error) {
	query := `
		SELECT id, user_id, subscription_id, amount, currency, payment_method, status, external_id, created_at, updated_at
		FROM payments
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(&payment.ID, &payment.UserID, &payment.SubscriptionID, &payment.Amount, &payment.Currency,
			&payment.PaymentMethod, &payment.Status, &payment.ExternalID, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payments: %w", err)
	}

	return payments, nil
}

// ListByStatus retrieves payments by status
func (r *PaymentRepository) ListByStatus(ctx context.Context, status models.PaymentStatus, limit, offset int) ([]*models.Payment, error) {
	query := `
		SELECT id, user_id, subscription_id, amount, currency, payment_method, status, external_id, created_at, updated_at
		FROM payments
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(&payment.ID, &payment.UserID, &payment.SubscriptionID, &payment.Amount, &payment.Currency,
			&payment.PaymentMethod, &payment.Status, &payment.ExternalID, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payments: %w", err)
	}

	return payments, nil
}

// Delete deletes a payment
func (r *PaymentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM payments WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// CountByUserID returns the total number of payments for a user
func (r *PaymentRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM payments WHERE user_id = $1`

	var count int64
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count payments: %w", err)
	}

	return count, nil
}
