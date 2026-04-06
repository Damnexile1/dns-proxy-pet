package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	DNS      DNSConfig      `yaml:"dns"`
	Proxy    ProxyConfig    `yaml:"proxy"`
	API      APIConfig      `yaml:"api"`
	Telegram TelegramConfig `yaml:"telegram"`
	Payment  PaymentConfig  `yaml:"payment"`
}

type AppConfig struct {
	Env      string `yaml:"env"`
	LogLevel string `yaml:"log_level"`
}

type DatabaseConfig struct {
	Host               string        `yaml:"host"`
	Port               int           `yaml:"port"`
	User               string        `yaml:"user"`
	Password           string        `yaml:"password"`
	Name               string        `yaml:"name"`
	SSLMode            string        `yaml:"sslmode"`
	MaxConnections     int           `yaml:"max_connections"`
	MaxIdleConnections int           `yaml:"max_idle_connections"`
	ConnectionLifetime time.Duration `yaml:"connection_lifetime"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

type DNSConfig struct {
	Port            int           `yaml:"port"`
	Host            string        `yaml:"host"`
	UpstreamServers []string      `yaml:"upstream_servers"`
	CacheTTL        int           `yaml:"cache_ttl"`
	Timeout         time.Duration `yaml:"timeout"`
}

type ProxyConfig struct {
	HTTPPort              int           `yaml:"http_port"`
	HTTPSPort             int           `yaml:"https_port"`
	Host                  string        `yaml:"host"`
	Timeout               time.Duration `yaml:"timeout"`
	MaxIdleConnections    int           `yaml:"max_idle_connections"`
	IdleConnectionTimeout time.Duration `yaml:"idle_connection_timeout"`
}

type APIConfig struct {
	Port          int           `yaml:"port"`
	Host          string        `yaml:"host"`
	ReadTimeout   time.Duration `yaml:"read_timeout"`
	WriteTimeout  time.Duration `yaml:"write_timeout"`
	JWTSecret     string        `yaml:"jwt_secret"`
	JWTExpiration time.Duration `yaml:"jwt_expiration"`
}

type TelegramConfig struct {
	BotToken   string  `yaml:"bot_token"`
	AdminIDs   []int64 `yaml:"admin_ids"`
	WebhookURL string  `yaml:"webhook_url"`
}

type PaymentConfig struct {
	YooKassa YooKassaConfig `yaml:"yookassa"`
}

type YooKassaConfig struct {
	ShopID        string `yaml:"shop_id"`
	SecretKey     string `yaml:"secret_key"`
	WebhookSecret string `yaml:"webhook_secret"`
	ReturnURL     string `yaml:"return_url"`
}

func Load(configPath string) (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables if set
	overrideFromEnv(&cfg)

	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	if v := os.Getenv("APP_ENV"); v != "" {
		cfg.App.Env = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.App.LogLevel = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.Name = v
	}
	if v := os.Getenv("REDIS_HOST"); v != "" {
		cfg.Redis.Host = v
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.API.JWTSecret = v
	}
	if v := os.Getenv("TELEGRAM_BOT_TOKEN"); v != "" {
		cfg.Telegram.BotToken = v
	}
	if v := os.Getenv("YOOKASSA_SHOP_ID"); v != "" {
		cfg.Payment.YooKassa.ShopID = v
	}
	if v := os.Getenv("YOOKASSA_SECRET_KEY"); v != "" {
		cfg.Payment.YooKassa.SecretKey = v
	}
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}
