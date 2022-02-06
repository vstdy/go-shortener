package psql

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	inter "github.com/vstdy0/go-project/storage"
	"github.com/vstdy0/go-project/storage/psql/schema"
)

const (
	serviceName = "psql"

	dbTableLoggingKey     = "db-table"
	dbOperationLoggingKey = "db-operation"
)

var _ inter.URLStorage = (*Storage)(nil)

type (
	Storage struct {
		config Config
		db     *bun.DB
	}

	StorageOption func(st *Storage) error
)

// WithConfig overrides default Storage config.
func WithConfig(config Config) StorageOption {
	return func(st *Storage) error {
		st.config = config

		return nil
	}
}

// New creates a new psql storage with custom options.
func New(opts ...StorageOption) (*Storage, error) {
	st := &Storage{
		config: NewDefaultConfig(),
	}
	for optIdx, opt := range opts {
		if err := opt(st); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	if err := st.config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(st.config.DSN)))

	maxOpenConnections := 4 * runtime.GOMAXPROCS(0)

	st.db = bun.NewDB(sqlDB, pgdialect.New())
	st.db.AddQueryHook(newQueryHook(st))
	st.db.SetMaxOpenConns(maxOpenConnections)
	st.db.SetMaxIdleConns(maxOpenConnections)
	st.db.RegisterModel(
		(*schema.URL)(nil),
	)

	if err := st.db.Ping(); err != nil {
		return nil, fmt.Errorf("ping for DSN (%s) failed: %w", st.config.DSN, err)
	}

	_, err := st.db.NewCreateTable().
		Model(&schema.URL{}).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return st, fmt.Errorf("create table failed: %w", err)
	}

	return st, nil
}

// Close closes DB connection.
func (st Storage) Close() error {
	if st.db == nil {
		return nil
	}

	return st.db.Close()
}

func (st *Storage) GetPing() error {
	if err := st.db.Ping(); err != nil {
		return err
	}

	return nil
}
