package rest

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

// Config keeps api params.
type Config struct {
	Timeout       time.Duration `mapstructure:"-"`
	LogLevel      zerolog.Level `mapstructure:"-"`
	ServerAddress string        `mapstructure:"server_address"`
	BaseURL       string        `mapstructure:"base_url"`
	SecretKey     string        `mapstructure:"secret_key"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.ServerAddress == "" {
		return fmt.Errorf("%s field: empty", "server_address")
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
		Timeout:       5 * time.Second,
		LogLevel:      zerolog.InfoLevel,
		ServerAddress: "0.0.0.0:8080",
		BaseURL:       "http://127.0.0.1:8080",
		SecretKey:     "secret_key",
	}
}
