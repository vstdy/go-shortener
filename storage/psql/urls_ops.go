package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"

	"github.com/vstdy0/go-shortener/model"
	"github.com/vstdy0/go-shortener/pkg"
	"github.com/vstdy0/go-shortener/storage/psql/schema"
)

const tableName = "url"

// HasURL checks existence of the url object with given id
func (st *Storage) HasURL(ctx context.Context, id int) (bool, error) {
	exists, err := st.db.NewSelect().
		Model(&schema.URL{}).
		Where("id = ?", id).
		Exists(ctx)

	return exists, err
}

// AddURLs adds given url objects to storage
func (st *Storage) AddURLs(ctx context.Context, objs []model.URL) ([]model.URL, error) {
	logger := st.Logger(withTable(tableName), withOperation("insert"))

	dbObjs := schema.NewURLsFromCanonical(objs)

	_, err := st.db.NewInsert().
		Model(&dbObjs).
		On("CONFLICT (url) WHERE deleted_at IS NULL DO UPDATE").
		Set("updated_at=NOW()").
		Returning("*, created_at <> updated_at AS updated").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	addedObjs, err := dbObjs.ToCanonical()
	if err != nil {
		return nil, err
	}

	logger.Info().Msgf("Objects added %v", addedObjs)

	for _, obj := range dbObjs {
		if obj.Updated {
			return addedObjs, pkg.ErrAlreadyExists
		}
	}

	return addedObjs, nil
}

// GetURL gets url object with given id
func (st *Storage) GetURL(ctx context.Context, id int) (model.URL, error) {
	dbObj := schema.URL{}

	err := st.db.NewSelect().
		Model(&dbObj).
		Where("id = ?", id).
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

// GetUserURLs gets current user url objects
func (st *Storage) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error) {
	var dbObjs schema.URLS

	err := st.db.NewSelect().
		Model(&dbObjs).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	if dbObjs == nil {
		return nil, nil
	}

	objs, err := dbObjs.ToCanonical()
	if err != nil {
		return nil, err
	}

	return objs, nil
}

// RemoveUserURLs removes current user url objects with given ids
func (st *Storage) RemoveUserURLs(ctx context.Context, objs []model.URL) error {
	logger := st.Logger(withTable(tableName), withOperation("delete"))

	dbObjs := schema.NewURLsFromCanonical(objs)

	res, err := st.db.NewDelete().
		Model(&dbObjs).
		WherePK("id", "user_id").
		Exec(ctx)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	logger.Info().Msgf("%d objects deleted", affected)

	return nil
}
