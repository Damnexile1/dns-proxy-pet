package models

import "time"

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod represents the payment method used
type PaymentMethod string

const (
	PaymentMethodYooKassa PaymentMethod = "yookassa"
	PaymentMethodCrypto   PaymentMethod = "crypto"
)

// Payment represents a payment transaction
type Payment struct {
	ID             int64         `json:"id" db:"id"`
	UserID         int64         `json:"user_id" db:"user_id"`
	SubscriptionID *int64        `json:"subscription_id" db:"subscription_id"`
	Amount         float64       `json:"amount" db:"amount"`
	Currency       string        `json:"currency" db:"currency"`
	PaymentMethod  PaymentMethod `json:"payment_method" db:"payment_method"`
	Status         PaymentStatus `json:"status" db:"status"`
	ExternalID     string        `json:"external_id" db:"external_id"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for Payment model
func (Payment) TableName() string {
	return "payments"
}

// IsCompleted checks if payment is completed
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusCompleted
}

// IsPending checks if payment is pending
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending
}
