package repository

import (
	"context"
	"fmt"

	"github.com/damnexile/dns-proxy-pet/internal/database"
	"github.com/damnexile/dns-proxy-pet/internal/models"
	"github.com/jackc/pgx/v5"
)

// ProxyCredentialsRepository handles database operations for proxy credentials
type ProxyCredentialsRepository struct {
	db database.Database
}

// NewProxyCredentialsRepository creates a new proxy credentials repository
func NewProxyCredentialsRepository(db database.Database) *ProxyCredentialsRepository {
	return &ProxyCredentialsRepository{db: db}
}

// Create creates new proxy credentials
func (r *ProxyCredentialsRepository) Create(ctx context.Context, creds *models.ProxyCredentials) error {
	query := `
		INSERT INTO proxy_credentials (user_id, token, login, password, max_devices, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		creds.UserID, creds.Token, creds.Login, creds.Password, creds.MaxDevices).
		Scan(&creds.ID, &creds.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create proxy credentials: %w", err)
	}

	return nil
}

// GetByID retrieves proxy credentials by ID
func (r *ProxyCredentialsRepository) GetByID(ctx context.Context, id int64) (*models.ProxyCredentials, error) {
	query := `
		SELECT id, user_id, token, login, password, max_devices, created_at
		FROM proxy_credentials
		WHERE id = $1
	`

	var creds models.ProxyCredentials
	err := r.db.QueryRow(ctx, query, id).
		Scan(&creds.ID, &creds.UserID, &creds.Token, &creds.Login, &creds.Password, &creds.MaxDevices, &creds.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("proxy credentials not found")
		}
		return nil, fmt.Errorf("failed to get proxy credentials: %w", err)
	}

	return &creds, nil
}

// GetByUserID retrieves proxy credentials by user ID
func (r *ProxyCredentialsRepository) GetByUserID(ctx context.Context, userID int64) (*models.ProxyCredentials, error) {
	query := `
		SELECT id, user_id, token, login, password, max_devices, created_at
		FROM proxy_credentials
		WHERE user_id = $1
	`

	var creds models.ProxyCredentials
	err := r.db.QueryRow(ctx, query, userID).
		Scan(&creds.ID, &creds.UserID, &creds.Token, &creds.Login, &creds.Password, &creds.MaxDevices, &creds.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("proxy credentials not found")
		}
		return nil, fmt.Errorf("failed to get proxy credentials: %w", err)
	}

	return &creds, nil
}

// GetByToken retrieves proxy credentials by token
func (r *ProxyCredentialsRepository) GetByToken(ctx context.Context, token string) (*models.ProxyCredentials, error) {
	query := `
		SELECT id, user_id, token, login, password, max_devices, created_at
		FROM proxy_credentials
		WHERE token = $1
	`

	var creds models.ProxyCredentials
	err := r.db.QueryRow(ctx, query, token).
		Scan(&creds.ID, &creds.UserID, &creds.Token, &creds.Login, &creds.Password, &creds.MaxDevices, &creds.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("proxy credentials not found")
		}
		return nil, fmt.Errorf("failed to get proxy credentials: %w", err)
	}

	return &creds, nil
}

// GetByLogin retrieves proxy credentials by login
func (r *ProxyCredentialsRepository) GetByLogin(ctx context.Context, login string) (*models.ProxyCredentials, error) {
	query := `
		SELECT id, user_id, token, login, password, max_devices, created_at
		FROM proxy_credentials
		WHERE login = $1
	`

	var creds models.ProxyCredentials
	err := r.db.QueryRow(ctx, query, login).
		Scan(&creds.ID, &creds.UserID, &creds.Token, &creds.Login, &creds.Password, &creds.MaxDevices, &creds.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("proxy credentials not found")
		}
		return nil, fmt.Errorf("failed to get proxy credentials: %w", err)
	}

	return &creds, nil
}

// Update updates proxy credentials
func (r *ProxyCredentialsRepository) Update(ctx context.Context, creds *models.ProxyCredentials) error {
	query := `
		UPDATE proxy_credentials
		SET token = $1, login = $2, password = $3, max_devices = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(ctx, query, creds.Token, creds.Login, creds.Password, creds.MaxDevices, creds.ID)
	if err != nil {
		return fmt.Errorf("failed to update proxy credentials: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("proxy credentials not found")
	}

	return nil
}

// Delete deletes proxy credentials
func (r *ProxyCredentialsRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM proxy_credentials WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete proxy credentials: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("proxy credentials not found")
	}

	return nil
}
