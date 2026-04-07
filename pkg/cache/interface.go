package cache

import (
	"context"
	"time"
)

// Cache is the interface for cache operations
// This abstraction allows us to switch cache implementations without changing business logic
type Cache interface {
	// Set sets a key-value pair with expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get gets a value by key
	Get(ctx context.Context, key string) (string, error)

	// Delete deletes a key
	Delete(ctx context.Context, keys ...string) error

	// Exists checks if a key exists
	Exists(ctx context.Context, keys ...string) (int64, error)

	// Expire sets expiration on a key
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// Increment increments a key by 1
	Increment(ctx context.Context, key string) (int64, error)

	// IncrementBy increments a key by a specific value
	IncrementBy(ctx context.Context, key string, value int64) (int64, error)

	// SetNX sets a key-value pair only if the key does not exist
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)

	// GetSet sets a new value and returns the old value
	GetSet(ctx context.Context, key string, value interface{}) (string, error)

	// FlushAll flushes all keys from the current database
	FlushAll(ctx context.Context) error

	// Ping checks if the cache connection is alive
	Ping(ctx context.Context) error

	// Close closes the cache connection
	Close() error
}
