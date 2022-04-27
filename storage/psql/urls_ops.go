package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/vstdy/go-shortener/model"
	"github.com/vstdy/go-shortener/pkg"
	"github.com/vstdy/go-shortener/pkg/tracing"
	"github.com/vstdy/go-shortener/storage/psql/schema"
)

const tableName = "url"

// HasURL checks existence of the url object with given id
func (st *Storage) HasURL(ctx context.Context, id int) (exists bool, err error) {
	ctx, span := tracing.StartSpanFromCtx(ctx, "psql HasURL")
	defer tracing.FinishSpan(span, err)

	logger := st.Logger(ctx, withTable(tableName), withOperation("HasURL"))

	exists, err = st.db.NewSelect().
		Model(&schema.URL{}).
		Where("id = ?", id).
		Exists(ctx)
	if err != nil {
		logger.Warn().Err(err).Msgf("check existence of URL with id: %v", id)
		return false, fmt.Errorf("psql: HasURL: %w", err)
	}

	return exists, err
}

// AddURLs adds given url objects to storage
func (st *Storage) AddURLs(ctx context.Context, objs []model.URL) (retObjs []model.URL, err error) {
	ctx, span := tracing.StartSpanFromCtx(ctx, "psql AddURLs")
	defer tracing.FinishSpan(span, err)

	logger := st.Logger(ctx, withTable(tableName), withOperation("AddURLs"))

	dbObjs := schema.NewURLsFromCanonical(objs)

	_, err = st.db.NewInsert().
		Model(&dbObjs).
		On("CONFLICT (url) WHERE deleted_at IS NULL DO UPDATE").
		Set("updated_at=NOW()").
		Returning("*, created_at <> updated_at AS updated").
		Exec(ctx)
	if err != nil {
		logger.Warn().Err(err).Msgf("add URLs: %v", dbObjs)
		return nil, fmt.Errorf("psql: AddURLs: %w", err)
	}

	retObjs, err = dbObjs.ToCanonical()
	if err != nil {
		return nil, fmt.Errorf("psql: AddURLs: converting to canonical: %w", err)
	}

	for _, obj := range dbObjs {
		if obj.Updated {
			return retObjs, fmt.Errorf("psql: AddURLs: %w", pkg.ErrAlreadyExists)
		}
	}

	return retObjs, nil
}

// GetURL gets url object with given id
func (st *Storage) GetURL(ctx context.Context, id int) (obj model.URL, err error) {
	ctx, span := tracing.StartSpanFromCtx(ctx, "psql GetURL")
	defer tracing.FinishSpan(span, err)

	logger := st.Logger(ctx, withTable(tableName), withOperation("GetURL"))

	dbObj := schema.URL{}

	err = st.db.NewSelect().
		Model(&dbObj).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.URL{}, nil
		}

		logger.Warn().Err(err).Msgf("get URL with id: %v", id)
		return model.URL{}, fmt.Errorf("psql: GetURL: %w", err)
	}

	obj = dbObj.ToCanonical()

	return obj, nil
}

// GetUsersURLs gets current user url objects
func (st *Storage) GetUsersURLs(ctx context.Context, userID uuid.UUID) (objs []model.URL, err error) {
	ctx, span := tracing.StartSpanFromCtx(ctx, "psql GetUsersURLs")
	defer tracing.FinishSpan(span, err)

	logger := st.Logger(ctx, withTable(tableName), withOperation("GetUsersURLs"))

	var dbObjs schema.URLS

	err = st.db.NewSelect().
		Model(&dbObjs).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		logger.Warn().Err(err).Msgf("get URLs of user with id: %v", userID)
		return nil, fmt.Errorf("psql: GetUsersURLs: %w", err)
	}
	if dbObjs == nil {
		return nil, nil
	}

	objs, err = dbObjs.ToCanonical()
	if err != nil {
		return nil, fmt.Errorf("psql: GetUsersURLs: converting to canonical: %w", err)
	}

	return objs, nil
}

// RemoveUsersURLs removes current user url objects with given ids
func (st *Storage) RemoveUsersURLs(ctx context.Context, objs []model.URL) (err error) {
	ctx, span := tracing.StartSpanFromCtx(ctx, "psql RemoveUsersURLs")
	defer tracing.FinishSpan(span, err)

	logger := st.Logger(ctx, withTable(tableName), withOperation("RemoveUsersURLs"))

	dbObjs := schema.NewURLsFromCanonical(objs)

	_, err = st.db.NewDelete().
		Model(&dbObjs).
		WherePK("id", "user_id").
		Exec(ctx)
	if err != nil {
		logger.Warn().Err(err).Msgf("remove following objects: %v", dbObjs)
		return fmt.Errorf("psql: RemoveUsersURLs: %w", err)
	}

	return nil
}
