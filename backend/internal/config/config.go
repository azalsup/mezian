// Package config loads configuration from config/config.yaml.
// Environment variables JWT_SECRET and DB_PATH override the file.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server       ServerConfig       `yaml:"server"`
	Database     DatabaseConfig     `yaml:"database"`
	JWT          JWTConfig          `yaml:"jwt"`
	OTP          OTPConfig          `yaml:"otp"`
	Notification NotificationConfig `yaml:"notification"`
	Media        MediaConfig        `yaml:"media"`
	Plans        PlansConfig        `yaml:"plans"`
}

type ServerConfig struct {
	Port        int      `yaml:"port"`
	Mode        string   `yaml:"mode"`
	BaseURL     string   `yaml:"base_url"`
	FrontendURL string   `yaml:"frontend_url"`
	CORSOrigins []string `yaml:"cors_origins"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type JWTConfig struct {
	Secret           string `yaml:"secret"`
	AccessTTLMinutes int    `yaml:"access_ttl_minutes"`
	RefreshTTLDays   int    `yaml:"refresh_ttl_days"`
}

type OTPConfig struct {
	Length           int `yaml:"length"`
	TTLMinutes       int `yaml:"ttl_minutes"`
	MaxAttempts      int `yaml:"max_attempts"`
	RateLimitPerHour int `yaml:"rate_limit_per_hour"`
}

type NotificationConfig struct {
	Provider string        `yaml:"provider"`
	Twilio   TwilioConfig  `yaml:"twilio"`
	WhatsApp WhatsAppConfig `yaml:"whatsapp"`
}

type TwilioConfig struct {
	AccountSID string `yaml:"account_sid"`
	AuthToken  string `yaml:"auth_token"`
	FromNumber string `yaml:"from_number"`
}

type WhatsAppConfig struct {
	APIURL        string `yaml:"api_url"`
	APIToken      string `yaml:"api_token"`
	PhoneNumberID string `yaml:"phone_number_id"`
}

type MediaConfig struct {
	UploadDir       string   `yaml:"upload_dir"`
	MaxSizeMB       int      `yaml:"max_size_mb"`
	MaxPerAd        int      `yaml:"max_per_ad"`
	ThumbnailWidth  int      `yaml:"thumbnail_width"`
	ThumbnailHeight int      `yaml:"thumbnail_height"`
	AllowedTypes    []string `yaml:"allowed_types"`
}

type PlansConfig struct {
	Starter PlanConfig `yaml:"starter"`
	Pro     PlanConfig `yaml:"pro"`
	Premium PlanConfig `yaml:"premium"`
}

type PlanConfig struct {
	PriceMAD     float64 `yaml:"price_mad"`
	MaxAds       int     `yaml:"max_ads"` // -1 = illimité
	DurationDays int     `yaml:"duration_days"`
}

// Load looks for config/config.local.yaml then config/config.yaml.
func Load() (*Config, error) {
	paths := []string{"config/config.local.yaml", "config/config.yaml"}
	var data []byte
	var err error
	for _, p := range paths {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("config file not found: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	// Surcharges via variables d'environnement
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("DB_PATH"); v != "" {
		cfg.Database.Path = v
	}
	if v := os.Getenv("PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Server.Port)
	}
	return &cfg, nil
}
