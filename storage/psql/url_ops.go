package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

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
func (st *Storage) AddURLS(ctx context.Context, urls []model.URL) ([]model.URL, error) {
	dbObjs := schema.NewURLsFromCanonical(urls)

	_, err := st.db.NewInsert().
		Model(&dbObjs).
		On("CONFLICT (url) DO UPDATE").
		Set("updated_at=NOW()").
		Returning("*, created_at <> updated_at AS updated").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	objs, err := dbObjs.ToCanonical()
	if err != nil {
		return nil, err
	}

	for _, obj := range dbObjs {
		if obj.Updated {
			return objs, pkg.ErrIntegrityViolation
		}
	}

	return objs, nil
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
