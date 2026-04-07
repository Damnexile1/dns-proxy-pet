package models

import "time"

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
)

// PlanType represents the type of subscription plan
type PlanType string

const (
	PlanTypeBasic   PlanType = "basic"
	PlanTypePremium PlanType = "premium"
	PlanTypeFamily  PlanType = "family"
)

// Subscription represents a user subscription
type Subscription struct {
	ID        int64              `json:"id" db:"id"`
	UserID    int64              `json:"user_id" db:"user_id"`
	PlanType  PlanType           `json:"plan_type" db:"plan_type"`
	Status    SubscriptionStatus `json:"status" db:"status"`
	StartedAt time.Time          `json:"started_at" db:"started_at"`
	ExpiresAt time.Time          `json:"expires_at" db:"expires_at"`
	AutoRenew bool               `json:"auto_renew" db:"auto_renew"`
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for Subscription model
func (Subscription) TableName() string {
	return "subscriptions"
}

// IsActive checks if subscription is currently active
func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive && s.ExpiresAt.After(time.Now())
}

// IsExpired checks if subscription has expired
func (s *Subscription) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
