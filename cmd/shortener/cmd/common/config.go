package common

import (
	"fmt"
	"time"

	"github.com/vstdy0/go-project/service/shortener/v1"
	"github.com/vstdy0/go-project/storage"
	"github.com/vstdy0/go-project/storage/file"
	"github.com/vstdy0/go-project/storage/memory"
	"github.com/vstdy0/go-project/storage/psql"
)

// Config combines sub-configs for all services, storages and providers.
type Config struct {
	RequestTimeout time.Duration
	ServerAddress  string      `mapstructure:"server_address"`
	BaseURL        string      `mapstructure:"base_url"`
	SecretKey      string      `mapstructure:"secret_key"`
	FileStorage    file.Config `mapstructure:"file_storage,squash"`
	PSQLStorage    psql.Config `mapstructure:"psql_storage,squash"`
}

const (
	fileStorage = "file"
	psqlStorage = "psql"
)

// BuildMemoryStorage builds memory.Storage dependency.
func (c Config) BuildMemoryStorage() (*memory.Storage, error) {
	st, err := memory.New()
	if err != nil {
		return nil, fmt.Errorf("building memory storage: %w", err)
	}

	return st, nil
}

// BuildFileStorage builds file.Storage dependency.
func (c Config) BuildFileStorage() (*file.Storage, error) {
	st, err := file.New(
		file.WithConfig(c.FileStorage),
	)
	if err != nil {
		return nil, fmt.Errorf("building file storage: %w", err)
	}

	return st, nil
}

// BuildPsqlStorage builds psql.Storage dependency.
func (c Config) BuildPsqlStorage() (*psql.Storage, error) {
	st, err := psql.New(
		psql.WithConfig(c.PSQLStorage),
	)
	if err != nil {
		return nil, fmt.Errorf("building psql storage: %w", err)
	}

	return st, nil
}

// BuildService builds shortener.Service dependency.
func (c Config) BuildService(storageType string) (*shortener.Service, error) {
	var st storage.URLStorage
	var err error

	switch storageType {
	case fileStorage:
		st, err = c.BuildFileStorage()
	case psqlStorage:
		st, err = c.BuildPsqlStorage()
	default:
		st, err = c.BuildMemoryStorage()
	}
	if err != nil {
		return nil, err
	}

	svc, err := shortener.NewService(shortener.WithStorage(st))
	if err != nil {
		return nil, fmt.Errorf("building service: %w", err)
	}

	return svc, nil
}
