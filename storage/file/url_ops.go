package file

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/storage/file/schema"
)

// HasURL checks existence of the object with given id
func (st *Storage) HasURL(ctx context.Context, urlID int) (bool, error) {
	st.RLock()
	defer st.RUnlock()

	_, ok := st.urls[urlID]

	return ok, nil
}

// AddURLS adds given objects to storage
func (st *Storage) AddURLS(ctx context.Context, objs []model.URL) ([]model.URL, error) {
	st.Lock()
	defer st.Unlock()

	dbObjs := schema.NewURLsFromCanonical(objs)

	for idx := range dbObjs {
		dbObjs[idx].ID = st.id

		if err := st.encoder.Encode(dbObjs[idx]); err != nil {
			return nil, err
		}

		st.urls[dbObjs[idx].ID] = dbObjs[idx]
		st.id++
	}

	addedObjs := dbObjs.ToCanonical()

	return addedObjs, nil
}

// GetURL gets object with given id
func (st *Storage) GetURL(ctx context.Context, urlID int) (model.URL, error) {
	st.RLock()
	defer st.RUnlock()

	url, ok := st.urls[urlID]
	if !ok {
		return model.URL{}, fmt.Errorf("url does not exist")
	}

	return url.ToCanonical(), nil
}

// GetUserURLs gets current user objects
func (st *Storage) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error) {
	st.RLock()
	defer st.RUnlock()

	var urls schema.URLS
	for _, v := range st.urls {
		if v.UserID == userID {
			urls = append(urls, v)
		}
	}

	return urls.ToCanonical(), nil
}

// RemoveUserURLs removes current user objects with given ids
func (st *Storage) RemoveUserURLs(objs []model.URL) error {

	return nil
}
