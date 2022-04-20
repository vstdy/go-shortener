package grpc

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

// Config keeps grpc params.
type Config struct {
	KapMaxConnectionIdle time.Duration `mapstructure:"kap_max_connection_idle"`
	KapMaxConnectionAge  time.Duration `mapstructure:"kap_max_connection_age"`
	KapTime              time.Duration `mapstructure:"kap_time"`
	KapTimeout           time.Duration `mapstructure:"kap_timeout"`
	Timeout              time.Duration `mapstructure:"timeout"`
	LogLevel             zerolog.Level `mapstructure:"-"`
	ServerAddress        string        `mapstructure:"grpc_server_address"`
	BaseURL              string        `mapstructure:"base_url"`
	SecretKey            string        `mapstructure:"secret_key"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.ServerAddress == "" {
		return fmt.Errorf("%s field: empty", "grpc_server_address")
	}

	if config.BaseURL == "" {
		return fmt.Errorf("%s field: empty", "base_url")
	}

	if config.SecretKey == "" {
		return fmt.Errorf("%s field: empty", "secret_key")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		KapMaxConnectionIdle: 5 * time.Minute,
		KapMaxConnectionAge:  5 * time.Minute,
		KapTime:              1 * time.Hour,
		KapTimeout:           5 * time.Second,
		Timeout:              5 * time.Second,
		LogLevel:             zerolog.InfoLevel,
		ServerAddress:        "0.0.0.0:50051",
		BaseURL:              "http://127.0.0.1:8080",
		SecretKey:            "secret_key",
	}
}
