package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Database is the interface for database operations
// This abstraction allows us to switch database implementations without changing business logic
type Database interface {
	// Query executes a query that returns rows
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)

	// QueryRow executes a query that is expected to return at most one row
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row

	// Exec executes a query that doesn't return rows
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)

	// Ping verifies a connection to the database is still alive
	Ping(ctx context.Context) error

	// Close closes the database connection
	Close()

	// Begin starts a transaction
	Begin(ctx context.Context) (Tx, error)
}

// Rows is the interface for query result rows
type Rows interface {
	// Close closes the rows
	Close()

	// Err returns any error that occurred during iteration
	Err() error

	// Next prepares the next row for scanning
	Next() bool

	// Scan copies the columns from the current row into dest values
	Scan(dest ...interface{}) error
}

// Row is the interface for a single row result
type Row interface {
	// Scan copies the columns from the row into dest values
	Scan(dest ...interface{}) error
}

// Tx is the interface for database transactions
type Tx interface {
	// Commit commits the transaction
	Commit(ctx context.Context) error

	// Rollback aborts the transaction
	Rollback(ctx context.Context) error

	// Query executes a query within the transaction
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)

	// QueryRow executes a query within the transaction
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row

	// Exec executes a query within the transaction
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

// pgxRows wraps pgx.Rows to implement our Rows interface
type pgxRows struct {
	rows pgx.Rows
}

func (r *pgxRows) Close() {
	r.rows.Close()
}

func (r *pgxRows) Err() error {
	return r.rows.Err()
}

func (r *pgxRows) Next() bool {
	return r.rows.Next()
}

func (r *pgxRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

// pgxRow wraps pgx.Row to implement our Row interface
type pgxRow struct {
	row pgx.Row
}

func (r *pgxRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

// pgxTx wraps pgx.Tx to implement our Tx interface
type pgxTx struct {
	tx pgx.Tx
}

func (t *pgxTx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *pgxTx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *pgxTx) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	rows, err := t.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRows{rows: rows}, nil
}

func (t *pgxTx) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	return &pgxRow{row: t.tx.QueryRow(ctx, sql, args...)}
}

func (t *pgxTx) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return t.tx.Exec(ctx, sql, args...)
}
