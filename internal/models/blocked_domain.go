package models

import "time"

// BlockedDomain represents a domain in the blocklist
type BlockedDomain struct {
	ID      int64     `json:"id" db:"id"`
	Domain  string    `json:"domain" db:"domain"`
	AddedAt time.Time `json:"added_at" db:"added_at"`
}

// TableName returns the table name for BlockedDomain model
func (BlockedDomain) TableName() string {
	return "blocked_domains"
}

// IsWildcard checks if the domain is a wildcard pattern (e.g., *.example.com)
func (b *BlockedDomain) IsWildcard() bool {
	return len(b.Domain) > 0 && b.Domain[0] == '*'
}
