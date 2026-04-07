package models_test

import (
	"testing"
	"time"

	"github.com/damnexile/dns-proxy-pet/internal/models"
)

// BlockedDomain tests
func TestBlockedDomain_IsWildcard(t *testing.T) {
	tests := []struct {
		name   string
		domain *models.BlockedDomain
		want   bool
	}{
		{"wildcard domain", &models.BlockedDomain{Domain: "*.example.com"}, true},
		{"regular domain", &models.BlockedDomain{Domain: "example.com"}, false},
		{"empty domain", &models.BlockedDomain{Domain: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.domain.IsWildcard(); got != tt.want {
				t.Errorf("BlockedDomain.IsWildcard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockedDomain_TableName(t *testing.T) {
	domain := models.BlockedDomain{}
	if got := domain.TableName(); got != "blocked_domains" {
		t.Errorf("BlockedDomain.TableName() = %v, want blocked_domains", got)
	}
}

// Payment tests
func TestPayment_IsCompleted(t *testing.T) {
	tests := []struct {
		name    string
		payment *models.Payment
		want    bool
	}{
		{"completed", &models.Payment{Status: models.PaymentStatusCompleted}, true},
		{"pending", &models.Payment{Status: models.PaymentStatusPending}, false},
		{"failed", &models.Payment{Status: models.PaymentStatusFailed}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.payment.IsCompleted(); got != tt.want {
				t.Errorf("Payment.IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPayment_TableName(t *testing.T) {
	Payment := models.Payment{}
	if got := Payment.TableName(); got != "payments" {
		t.Errorf("Payment.TableName() = %v, want payments", got)
	}
}

// Subscription tests
func TestSubscription_IsActive(t *testing.T) {
	tests := []struct {
		name string
		sub  *models.Subscription
		want bool
	}{
		{"active future", &models.Subscription{Status: models.SubscriptionStatusActive, ExpiresAt: time.Now().Add(24 * time.Hour)}, true},
		{"active past", &models.Subscription{Status: models.SubscriptionStatusActive, ExpiresAt: time.Now().Add(-24 * time.Hour)}, false},
		{"expired", &models.Subscription{Status: models.SubscriptionStatusExpired, ExpiresAt: time.Now().Add(24 * time.Hour)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.IsActive(); got != tt.want {
				t.Errorf("Subscription.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscription_TableName(t *testing.T) {
	Subscription := models.Subscription{}
	if got := Subscription.TableName(); got != "subscriptions" {
		t.Errorf("Subscription.TableName() = %v, want subscriptions", got)
	}
}

// TrafficUsage tests
func TestTrafficUsage_GetGigabytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  float64
	}{
		{"1 GB", 1073741824, 1.0},
		{"5 GB", 5368709120, 5.0},
		{"zero", 0, 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traffic := &models.TrafficUsage{BytesUsed: tt.bytes}
			if got := traffic.GetGigabytes(); got != tt.want {
				t.Errorf("TrafficUsage.GetGigabytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrafficUsage_TableName(t *testing.T) {
	TrafficUsage := models.TrafficUsage{}
	if got := TrafficUsage.TableName(); got != "traffic_usage" {
		t.Errorf("TrafficUsage.TableName() = %v, want traffic_usage", got)
	}
}

// User tests
func TestUser_TableName(t *testing.T) {
	User := models.User{}
	if got := User.TableName(); got != "users" {
		t.Errorf("User.TableName() = %v, want users", got)
	}
}

// ProxyCredentials tests
func TestProxyCredentials_TableName(t *testing.T) {
	ProxyCredentials := models.ProxyCredentials{}
	if got := ProxyCredentials.TableName(); got != "proxy_credentials" {
		t.Errorf("ProxyCredentials.TableName() = %v, want proxy_credentials", got)
	}
}
