package models

import "time"

// User represents a user in the system
type User struct {
	ID         int64     `json:"id" db:"id"`
	TelegramID int64     `json:"telegram_id" db:"telegram_id"`
	Username   string    `json:"username" db:"username"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}
