package psql

import "fmt"

const (
	defaultConfigEndpoint = "postgres://username:password@localhost:5432/database_name"
)

// Config keeps Storage configuration.
type Config struct {
	DSN string `mapstructure:"database_dsn"`
}

// Validate performs a basic validation.
func (c Config) Validate() error {
	if c.DSN == "" {
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
