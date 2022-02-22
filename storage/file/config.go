package file

import (
	"fmt"
	"path/filepath"
	"runtime"
)

const (
	defaultFileStorageName = "storage_file.txt"
)

// Config keeps Storage configuration.
type Config struct {
	FileStoragePath string `mapstructure:"file_storage_path"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.FileStoragePath == "" {
		return fmt.Errorf("%s field: empty", "FileStoragePath")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		FileStoragePath: defaultFileStoragePath(),
	}
}

// defaultFileStoragePath returns full path to storage file.
func defaultFileStoragePath() string {
	_, filePath, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	return filepath.Dir(filePath) + "/" + defaultFileStorageName
}
