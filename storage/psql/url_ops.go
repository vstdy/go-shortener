package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/vstdy0/go-project/model"
	inter "github.com/vstdy0/go-project/storage"
	"github.com/vstdy0/go-project/storage/psql/schema"
)

var _ inter.URLStorage = (*Storage)(nil)

type (
	Storage struct {
		config Config
		db     *bun.DB
	}

	StorageOption func(st *Storage) error
)

// Has checks existence of the object with given id
func (st *Storage) Has(ctx context.Context, id int) (bool, error) {
	exists, err := st.db.NewSelect().
		Model(&schema.URL{}).
		Where("id = ?", id).
		Exists(ctx)

	return exists, err
}

func (st *Storage) Set(ctx context.Context, urls []model.URL) ([]model.URL, error) {
	dbObjs := schema.NewURLsFromCanonical(urls)

	_, err := st.db.NewInsert().
		Model(&dbObjs).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	objs, err := dbObjs.ToCanonical()
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (st *Storage) Get(ctx context.Context, id int) (model.URL, error) {
	dbObj := &schema.URL{}

	err := st.db.NewSelect().
		Model(dbObj).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.URL{}, nil
		}
		return model.URL{}, err
	}

	obj := dbObj.ToCanonical()

	return obj, nil
}

func (st *Storage) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error) {
	var dbObjs schema.URLS

	err := st.db.NewSelect().
		Model(&dbObjs).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	objs, err := dbObjs.ToCanonical()
	if err != nil {
		return nil, fmt.Errorf("conveting to canonical models: %w", err)
	}

	return objs, nil
}

// WithConfig overrides default Storage config.
func WithConfig(config Config) StorageOption {
	return func(st *Storage) error {
		st.config = config

		return nil
	}
}

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

	db := bun.NewDB(sqlDB, pgdialect.New())
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxOpenConnections)
	db.RegisterModel(
		(*schema.URL)(nil),
	)

	if err := db.Ping(); err != nil {
		return st, fmt.Errorf("ping for DSN (%s) failed: %w", st.config.DSN, err)
	}

	_, err := db.NewCreateTable().
		Model(&schema.URL{}).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return st, fmt.Errorf("create table failed: %w", err)
	}

	st.db = db

	return st, nil
}
