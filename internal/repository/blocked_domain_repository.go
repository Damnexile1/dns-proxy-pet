package repository

import (
	"context"
	"fmt"

	"github.com/damnexile/dns-proxy-pet/internal/database"
	"github.com/damnexile/dns-proxy-pet/internal/models"
	"github.com/jackc/pgx/v5"
)

// BlockedDomainRepository handles database operations for blocked domains
type BlockedDomainRepository struct {
	db database.Database
}

// NewBlockedDomainRepository creates a new blocked domain repository
func NewBlockedDomainRepository(db database.Database) *BlockedDomainRepository {
	return &BlockedDomainRepository{db: db}
}

// Create creates a new blocked domain
func (r *BlockedDomainRepository) Create(ctx context.Context, domain *models.BlockedDomain) error {
	query := `
		INSERT INTO blocked_domains (domain, added_at)
		VALUES ($1, NOW())
		RETURNING id, added_at
	`

	err := r.db.QueryRow(ctx, query, domain.Domain).
		Scan(&domain.ID, &domain.AddedAt)
	if err != nil {
		return fmt.Errorf("failed to create blocked domain: %w", err)
	}

	return nil
}

// GetByID retrieves a blocked domain by ID
func (r *BlockedDomainRepository) GetByID(ctx context.Context, id int64) (*models.BlockedDomain, error) {
	query := `
		SELECT id, domain, added_at
		FROM blocked_domains
		WHERE id = $1
	`

	var domain models.BlockedDomain
	err := r.db.QueryRow(ctx, query, id).
		Scan(&domain.ID, &domain.Domain, &domain.AddedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("blocked domain not found")
		}
		return nil, fmt.Errorf("failed to get blocked domain: %w", err)
	}

	return &domain, nil
}

// GetByDomain retrieves a blocked domain by domain name
func (r *BlockedDomainRepository) GetByDomain(ctx context.Context, domainName string) (*models.BlockedDomain, error) {
	query := `
		SELECT id, domain, added_at
		FROM blocked_domains
		WHERE domain = $1
	`

	var domain models.BlockedDomain
	err := r.db.QueryRow(ctx, query, domainName).
		Scan(&domain.ID, &domain.Domain, &domain.AddedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("blocked domain not found")
		}
		return nil, fmt.Errorf("failed to get blocked domain: %w", err)
	}

	return &domain, nil
}

// IsBlocked checks if a domain is blocked
func (r *BlockedDomainRepository) IsBlocked(ctx context.Context, domainName string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM blocked_domains WHERE domain = $1)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, domainName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if domain is blocked: %w", err)
	}

	return exists, nil
}

// List retrieves all blocked domains with pagination
func (r *BlockedDomainRepository) List(ctx context.Context, limit, offset int) ([]*models.BlockedDomain, error) {
	query := `
		SELECT id, domain, added_at
		FROM blocked_domains
		ORDER BY added_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list blocked domains: %w", err)
	}
	defer rows.Close()

	var domains []*models.BlockedDomain
	for rows.Next() {
		var domain models.BlockedDomain
		err := rows.Scan(&domain.ID, &domain.Domain, &domain.AddedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan blocked domain: %w", err)
		}
		domains = append(domains, &domain)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating blocked domains: %w", err)
	}

	return domains, nil
}

// ListAll retrieves all blocked domains without pagination
func (r *BlockedDomainRepository) ListAll(ctx context.Context) ([]*models.BlockedDomain, error) {
	query := `
		SELECT id, domain, added_at
		FROM blocked_domains
		ORDER BY domain ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all blocked domains: %w", err)
	}
	defer rows.Close()

	var domains []*models.BlockedDomain
	for rows.Next() {
		var domain models.BlockedDomain
		err := rows.Scan(&domain.ID, &domain.Domain, &domain.AddedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan blocked domain: %w", err)
		}
		domains = append(domains, &domain)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating blocked domains: %w", err)
	}

	return domains, nil
}

// Delete deletes a blocked domain
func (r *BlockedDomainRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blocked_domains WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete blocked domain: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("blocked domain not found")
	}

	return nil
}

// DeleteByDomain deletes a blocked domain by domain name
func (r *BlockedDomainRepository) DeleteByDomain(ctx context.Context, domainName string) error {
	query := `DELETE FROM blocked_domains WHERE domain = $1`

	result, err := r.db.Exec(ctx, query, domainName)
	if err != nil {
		return fmt.Errorf("failed to delete blocked domain: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("blocked domain not found")
	}

	return nil
}

// Count returns the total number of blocked domains
func (r *BlockedDomainRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM blocked_domains`

	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count blocked domains: %w", err)
	}

	return count, nil
}
