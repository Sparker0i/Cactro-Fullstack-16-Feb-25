package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	RateLimit  RateLimitConfig
	Cors       CorsConfig
	Logger     LoggerConfig
	Monitoring MonitoringConfig
}

type ServerConfig struct {
	Port            string        `envconfig:"SERVER_PORT" default:"8080"`
	Host            string        `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	Mode            string        `envconfig:"SERVER_MODE" default:"release"`
	TimeoutRead     time.Duration `envconfig:"SERVER_TIMEOUT_READ" default:"5s"`
	TimeoutWrite    time.Duration `envconfig:"SERVER_TIMEOUT_WRITE" default:"10s"`
	TimeoutIdle     time.Duration `envconfig:"SERVER_TIMEOUT_IDLE" default:"120s"`
	ShutdownTimeout time.Duration `envconfig:"SERVER_SHUTDOWN_TIMEOUT" default:"20s"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" default:"polling_app"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
	MaxConns int32  `envconfig:"DB_MAX_CONNS" default:"25"`
	MinConns int32  `envconfig:"DB_MIN_CONNS" default:"5"`
}

type RateLimitConfig struct {
	Enabled           bool          `envconfig:"RATE_LIMIT_ENABLED" default:"true"`
	RequestsPerMinute int           `envconfig:"RATE_LIMIT_REQUESTS" default:"100"`
	BurstSize         int           `envconfig:"RATE_LIMIT_BURST" default:"20"`
	TTL               time.Duration `envconfig:"RATE_LIMIT_TTL" default:"1m"`
}

type CorsConfig struct {
	AllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
	AllowedMethods []string `envconfig:"CORS_ALLOWED_METHODS" default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders []string `envconfig:"CORS_ALLOWED_HEADERS" default:"Origin,Content-Type,Accept,Authorization"`
	MaxAge         int      `envconfig:"CORS_MAX_AGE" default:"300"`
}

type LoggerConfig struct {
	Level  string `envconfig:"LOG_LEVEL" default:"info"`
	Format string `envconfig:"LOG_FORMAT" default:"json"`
	Output string `envconfig:"LOG_OUTPUT" default:"stdout"`
}

type MonitoringConfig struct {
	Enabled     bool   `envconfig:"MONITORING_ENABLED" default:"true"`
	MetricsPort string `envconfig:"METRICS_PORT" default:"9090"`
}

func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Generate database connection string
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}
