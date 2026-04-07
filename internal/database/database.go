package database

import (
	"context"
	"fmt"
	"time"

	"github.com/damnexile/dns-proxy-pet/internal/config"
	"github.com/damnexile/dns-proxy-pet/pkg/logger"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// PgxAdapter is the PostgreSQL adapter that implements Database interface
type PgxAdapter struct {
	pool *pgxpool.Pool
}

// New creates a new database connection pool with PgxAdapter
func New(cfg *config.DatabaseConfig) (Database, error) {
	dsn := cfg.DSN()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set connection pool settings
	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MinConns = int32(cfg.MaxIdleConnections)
	poolConfig.MaxConnLifetime = cfg.ConnectionLifetime
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Name),
	)

	return &PgxAdapter{pool: pool}, nil
}

// Query executes a query that returns rows
func (db *PgxAdapter) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRows{rows: rows}, nil
}

// QueryRow executes a query that is expected to return at most one row
func (db *PgxAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	return &pgxRow{row: db.pool.QueryRow(ctx, sql, args...)}
}

// Exec executes a query that doesn't return rows
func (db *PgxAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, sql, args...)
}

// Ping checks if the database connection is alive
func (db *PgxAdapter) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Close closes the database connection pool
func (db *PgxAdapter) Close() {
	if db.pool != nil {
		db.pool.Close()
		logger.Info("Database connection closed")
	}
}

// Begin starts a transaction
func (db *PgxAdapter) Begin(ctx context.Context) (Tx, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &pgxTx{tx: tx}, nil
}

// Stats returns connection pool statistics
func (db *PgxAdapter) Stats() *pgxpool.Stat {
	return db.pool.Stat()
}
