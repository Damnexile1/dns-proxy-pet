package models

import "time"

// TrafficUsage represents daily traffic usage for a user
type TrafficUsage struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	BytesUsed int64     `json:"bytes_used" db:"bytes_used"`
	Date      time.Time `json:"date" db:"date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TableName returns the table name for TrafficUsage model
func (TrafficUsage) TableName() string {
	return "traffic_usage"
}

// GetMegabytes returns traffic usage in megabytes
func (t *TrafficUsage) GetMegabytes() float64 {
	return float64(t.BytesUsed) / 1024 / 1024
}

// GetGigabytes returns traffic usage in gigabytes
func (t *TrafficUsage) GetGigabytes() float64 {
	return float64(t.BytesUsed) / 1024 / 1024 / 1024
}

// AddBytes adds bytes to the current usage
func (t *TrafficUsage) AddBytes(bytes int64) {
	t.BytesUsed += bytes
}
