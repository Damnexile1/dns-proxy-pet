package models

import "time"

// ProxyCredentials represents proxy authentication credentials for a user
type ProxyCredentials struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	Token      string    `json:"token" db:"token"`
	Login      string    `json:"login" db:"login"`
	Password   string    `json:"password" db:"password"`
	MaxDevices int       `json:"max_devices" db:"max_devices"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// TableName returns the table name for ProxyCredentials model
func (ProxyCredentials) TableName() string {
	return "proxy_credentials"
}
