package common

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/vstdy/go-shortener/api/grpc"
	"github.com/vstdy/go-shortener/api/rest"
	"github.com/vstdy/go-shortener/pkg"
	"github.com/vstdy/go-shortener/service/shortener/v1"
	"github.com/vstdy/go-shortener/storage"
	"github.com/vstdy/go-shortener/storage/file"
	"github.com/vstdy/go-shortener/storage/memory"
	"github.com/vstdy/go-shortener/storage/psql"
)

// Config combines sub-configs for all services, storages and providers.
type Config struct {
	Timeout     time.Duration    `mapstructure:"timeout"`
	LogLevel    zerolog.Level    `mapstructure:"-"`
	StorageType string           `mapstructure:"storage_type"`
	HTTPServer  rest.Config      `mapstructure:"server,squash"`
	GRPCServer  grpc.Config      `mapstructure:"grpc_server,squash"`
	Service     shortener.Config `mapstructure:"service,squash"`
	FileStorage file.Config      `mapstructure:"file_storage,squash"`
	PSQLStorage psql.Config      `mapstructure:"psql_storage,squash"`
}

const (
	memoryStorage = "memory"
	fileStorage   = "file"
	psqlStorage   = "psql"
)

// BuildDefaultConfig builds a Config with default values.
func BuildDefaultConfig() Config {
	return Config{
		Timeout:     5 * time.Second,
		LogLevel:    zerolog.InfoLevel,
		StorageType: psqlStorage,
		HTTPServer:  rest.NewDefaultConfig(),
		GRPCServer:  grpc.NewDefaultConfig(),
		Service:     shortener.NewDefaultConfig(),
		FileStorage: file.NewDefaultConfig(),
		PSQLStorage: psql.NewDefaultConfig(),
	}
}

// BuildMemoryStorage builds memory.Storage dependency.
func (config Config) BuildMemoryStorage() (*memory.Storage, error) {
	st, err := memory.NewStorage()
	if err != nil {
		return nil, fmt.Errorf("building memory storage: %w", err)
	}

	return st, nil
}

// BuildFileStorage builds file.Storage dependency.
func (config Config) BuildFileStorage() (*file.Storage, error) {
	st, err := file.NewStorage(
		file.WithConfig(config.FileStorage),
	)
	if err != nil {
		return nil, fmt.Errorf("building file storage: %w", err)
	}

	return st, nil
}

// BuildPsqlStorage builds psql.Storage dependency.
func (config Config) BuildPsqlStorage() (*psql.Storage, error) {
	st, err := psql.NewStorage(
		psql.WithConfig(config.PSQLStorage),
	)
	if err != nil {
		return nil, fmt.Errorf("building psql storage: %w", err)
	}

	return st, nil
}

// BuildService builds shortener.Service dependency.
func (config Config) BuildService() (*shortener.Service, error) {
	var st storage.Storage
	var err error

	switch config.StorageType {
	case memoryStorage:
		st, err = config.BuildMemoryStorage()
	case fileStorage:
		st, err = config.BuildFileStorage()
	case psqlStorage:
		st, err = config.BuildPsqlStorage()
	default:
		err = pkg.ErrUnsupportedStorageType
	}
	if err != nil {
		return nil, err
	}

	svc, err := shortener.NewService(
		shortener.WithConfig(config.Service),
		shortener.WithStorage(st),
	)
	if err != nil {
		return nil, fmt.Errorf("building service: %w", err)
	}

	return svc, nil
}
