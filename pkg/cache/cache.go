package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/damnexile/dns-proxy-pet/internal/config"
	"github.com/damnexile/dns-proxy-pet/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisAdapter is the Redis adapter that implements Cache interface
type RedisAdapter struct {
	client *redis.Client
}

// New creates a new Redis cache client
func New(cfg *config.RedisConfig) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis connection established",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	return &RedisAdapter{client: client}, nil
}

// Close closes the Redis connection
func (c *RedisAdapter) Close() error {
	if c.client != nil {
		err := c.client.Close()
		logger.Info("Redis connection closed")
		return err
	}
	return nil
}

// Ping checks if the Redis connection is alive
func (c *RedisAdapter) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Set sets a key-value pair with expiration
func (c *RedisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (c *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Delete deletes a key
func (c *RedisAdapter) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (c *RedisAdapter) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

// Expire sets expiration on a key
func (c *RedisAdapter) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// Increment increments a key by 1
func (c *RedisAdapter) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// IncrementBy increments a key by a specific value
func (c *RedisAdapter) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, key, value).Result()
}

// SetNX sets a key-value pair only if the key does not exist
func (c *RedisAdapter) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, value, expiration).Result()
}

// GetSet sets a new value and returns the old value
func (c *RedisAdapter) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	return c.client.GetSet(ctx, key, value).Result()
}

// FlushAll flushes all keys from the current database
func (c *RedisAdapter) FlushAll(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}
