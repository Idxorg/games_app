package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

// Config holds all application configuration.
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Redis     RedisConfig     `yaml:"redis"`
	S3        S3Config        `yaml:"s3"`
	PortalAPI PortalAPIConfig `yaml:"portal_api"`
	JWT       JWTConfig       `yaml:"jwt"`
	Embed     EmbedConfig     `yaml:"embed"`
	LiveKit   LiveKitConfig   `yaml:"livekit"`
	CORS      CORSConfig      `yaml:"cors"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

type RateLimitConfig struct {
	MaxRequests   int `yaml:"max_requests"`
	WindowSeconds int `yaml:"window_seconds"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	URL      string `yaml:"url"`
	PoolSize int    `yaml:"pool_size"`
}

type RedisConfig struct {
	URL      string `yaml:"url"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type S3Config struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
}

type PortalAPIConfig struct {
	URL     string `yaml:"url"`
	APIKey  string `yaml:"api_key"`
	Timeout int    `yaml:"timeout"`
}

type JWTConfig struct {
	Secret       string `yaml:"secret"`
	JWKSURL      string `yaml:"jwks_url"`
	Issuer       string `yaml:"issuer"`
	ExpiryHours  int    `yaml:"expiry_hours"`
}

// EmbedConfig holds portal embed-specific configuration.
type EmbedConfig struct {
	HandoffSecret string `yaml:"handoff_secret"`
}

type LiveKitConfig struct {
	Host     string `yaml:"host"`
	APIKey   string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

// Load reads configuration from the given YAML file path and overlays
// environment variables. Environment variables take precedence over file values.
//
// Env var mapping:
//
//	SERVER_PORT, SERVER_HOST, SERVER_MODE
//	DATABASE_URL, DATABASE_POOL_SIZE
//	REDIS_URL, REDIS_PASSWORD, REDIS_DB
//	S3_ENDPOINT, S3_ACCESS_KEY, S3_SECRET_KEY, S3_BUCKET, S3_REGION
//	PORTAL_API_URL, PORTAL_API_KEY, PORTAL_API_TIMEOUT
//	JWT_SECRET, JWT_JWKS_URL, JWT_ISSUER
//	LIVEKIT_HOST, LIVEKIT_API_KEY, LIVEKIT_API_SECRET
//	CORS_ALLOWED_ORIGINS (comma-separated)
func Load(path string) (*Config, error) {
	cfg := &Config{}

	// Defaults
	cfg.Server.Port = 8000
	cfg.Server.Host = "0.0.0.0"
	cfg.Server.Mode = "release"
	cfg.Database.PoolSize = 10
	cfg.Redis.DB = 0
	cfg.PortalAPI.Timeout = 10
	cfg.S3.Region = "us-east-1"
	cfg.RateLimit.MaxRequests = 100
	cfg.RateLimit.WindowSeconds = 60
	cfg.JWT.ExpiryHours = 24

	// Load from file if it exists
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("reading config file %s: %w", path, err)
		}
		if len(data) > 0 {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("parsing config file %s: %w", path, err)
			}
		}
	}

	// Override with environment variables
	overrideString("SERVER_PORT", func(v string) { cfg.Server.Port = mustInt(v) })
	overrideString("SERVER_HOST", func(v string) { cfg.Server.Host = v })
	overrideString("SERVER_MODE", func(v string) { cfg.Server.Mode = v })
	overrideString("DATABASE_URL", func(v string) { cfg.Database.URL = v })
	overrideString("DATABASE_POOL_SIZE", func(v string) { cfg.Database.PoolSize = mustInt(v) })
	overrideString("REDIS_URL", func(v string) { cfg.Redis.URL = v })
	overrideString("REDIS_PASSWORD", func(v string) { cfg.Redis.Password = v })
	overrideString("REDIS_DB", func(v string) { cfg.Redis.DB = mustInt(v) })
	overrideString("S3_ENDPOINT", func(v string) { cfg.S3.Endpoint = v })
	overrideString("S3_ACCESS_KEY", func(v string) { cfg.S3.AccessKey = v })
	overrideString("S3_SECRET_KEY", func(v string) { cfg.S3.SecretKey = v })
	overrideString("S3_BUCKET", func(v string) { cfg.S3.Bucket = v })
	overrideString("S3_REGION", func(v string) { cfg.S3.Region = v })
	overrideString("PORTAL_API_URL", func(v string) { cfg.PortalAPI.URL = v })
	overrideString("PORTAL_API_KEY", func(v string) { cfg.PortalAPI.APIKey = v })
	overrideString("PORTAL_API_TIMEOUT", func(v string) { cfg.PortalAPI.Timeout = mustInt(v) })
	overrideString("JWT_SECRET", func(v string) { cfg.JWT.Secret = v })
	overrideString("JWT_JWKS_URL", func(v string) { cfg.JWT.JWKSURL = v })
	overrideString("JWT_ISSUER", func(v string) { cfg.JWT.Issuer = v })
	overrideString("JWT_EXPIRY_HOURS", func(v string) { cfg.JWT.ExpiryHours = mustInt(v) })
	overrideString("GAMES_EMBED_HANDOFF_SECRET", func(v string) { cfg.Embed.HandoffSecret = v })
	overrideString("LIVEKIT_HOST", func(v string) { cfg.LiveKit.Host = v })
	overrideString("LIVEKIT_API_KEY", func(v string) { cfg.LiveKit.APIKey = v })
	overrideString("LIVEKIT_API_SECRET", func(v string) { cfg.LiveKit.APISecret = v })
	overrideString("CORS_ALLOWED_ORIGINS", func(v string) {
		cfg.CORS.AllowedOrigins = strings.Split(v, ",")
		for i := range cfg.CORS.AllowedOrigins {
			cfg.CORS.AllowedOrigins[i] = strings.TrimSpace(cfg.CORS.AllowedOrigins[i])
		}
	})
	overrideString("RATE_LIMIT_MAX_REQUESTS", func(v string) { cfg.RateLimit.MaxRequests = mustInt(v) })
	overrideString("RATE_LIMIT_WINDOW_SECONDS", func(v string) { cfg.RateLimit.WindowSeconds = mustInt(v) })

	// Backwards compat: PORT env var maps to SERVER_PORT (used in old main.go)
	overrideString("PORT", func(v string) { cfg.Server.Port = mustInt(v) })

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that required configuration fields are set.
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	if c.Server.Host == "" {
		return fmt.Errorf("server host is required")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.Redis.URL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	return nil
}

// Addr returns the host:port string for the server.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// PortalAPITimeout returns the portal API timeout as a time.Duration.
func (c *Config) PortalAPITimeout() time.Duration {
	return time.Duration(c.PortalAPI.Timeout) * time.Second
}

func overrideString(key string, apply func(string)) {
	if v := os.Getenv(key); v != "" {
		apply(v)
	}
}

func mustInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}
