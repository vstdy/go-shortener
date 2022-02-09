package file

import "fmt"

const (
	defaultConfigEndpoint = "./storage/file/storage.txt"
)

// Config keeps FileStorage configuration.
type Config struct {
	FileStoragePath string `mapstructure:"file_storage_path"`
}

// Validate performs a basic validation.
func (c Config) Validate() error {
	if c.FileStoragePath == "" {
		return fmt.Errorf("%s field: empty", "FileStoragePath")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		FileStoragePath: defaultConfigEndpoint,
	}
}
