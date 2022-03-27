package memory

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vstdy0/go-shortener/model"
	"github.com/vstdy0/go-shortener/storage/memory/schema"
)

// HasURL checks existence of the url object with given id
func (st *Storage) HasURL(ctx context.Context, urlID int) (bool, error) {
	st.RLock()
	defer st.RUnlock()

	_, ok := st.urls[urlID]

	return ok, nil
}

// AddURLs adds given url objects to storage
func (st *Storage) AddURLs(ctx context.Context, objs []model.URL) ([]model.URL, error) {
	st.Lock()
	defer st.Unlock()

	dbObjs := schema.NewURLsFromCanonical(objs)

	for idx := range dbObjs {
		dbObjs[idx].ID = st.id

		st.urls[dbObjs[idx].ID] = dbObjs[idx]
		st.id++
	}

	addedObjs := dbObjs.ToCanonical()

	return addedObjs, nil
}

// GetURL gets url object with given id
func (st *Storage) GetURL(ctx context.Context, urlID int) (model.URL, error) {
	st.RLock()
	defer st.RUnlock()

	url, ok := st.urls[urlID]
	if !ok {
		return model.URL{}, fmt.Errorf("url does not exist")
	}

	return url.ToCanonical(), nil
}

// GetUserURLs gets current user url objects
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

// RemoveUserURLs removes current user url objects with given ids
func (st *Storage) RemoveUserURLs(ctx context.Context, objs []model.URL) error {

	return nil
}
