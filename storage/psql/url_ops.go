package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/pkg"
	"github.com/vstdy0/go-project/storage/psql/schema"
)

// HasURL checks existence of the object with given id
func (st *Storage) HasURL(ctx context.Context, id int) (bool, error) {
	exists, err := st.db.NewSelect().
		Model(&schema.URL{}).
		Where("id = ?", id).
		Exists(ctx)

	return exists, err
}

// AddURLS adds given objects to storage
func (st *Storage) AddURLS(ctx context.Context, urls []model.URL) ([]model.URL, error, error) {
	var pgErr error
	dbObjs := schema.NewURLsFromCanonical(urls)

	_, err := st.db.NewInsert().
		Model(&dbObjs).
		Returning("*").
		Exec(ctx)
	if err != nil {
		if err, ok := err.(pgdriver.Error); ok &&
			err.Field('C') == pgerrcode.UniqueViolation {
			_, err := st.db.NewInsert().
				Model(&dbObjs).
				On("CONFLICT (url) DO UPDATE").
				Returning("*").
				Exec(ctx)
			if err != nil {
				return nil, nil, err
			}

			pgErr = pkg.ErrUniqueViolation
		} else {
			return nil, nil, err
		}
	}

	objs, err := dbObjs.ToCanonical()
	if err != nil {
		return nil, nil, err
	}

	return objs, pgErr, nil
}

// GetURL gets object with given id
func (st *Storage) GetURL(ctx context.Context, id int) (model.URL, error) {
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

// GetUserURLs gets current user objects
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
