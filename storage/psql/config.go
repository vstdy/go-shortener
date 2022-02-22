package psql

import "fmt"

const (
	defaultConfigEndpoint = "postgres://user:password@localhost:5432/shortener?sslmode=disable"
)

// Config keeps Storage configuration.
type Config struct {
	DSN string `mapstructure:"database_dsn"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.DSN == "" {
		return fmt.Errorf("%s field: empty", "DSN")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		DSN: defaultConfigEndpoint,
	}
}
