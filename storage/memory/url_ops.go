package memory

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/storage/memory/schema"
)

// HasURL checks existence of the object with given id
func (st *Storage) HasURL(ctx context.Context, urlID int) (bool, error) {
	st.RLock()
	defer st.RUnlock()
	_, ok := st.urls[urlID]

	return ok, nil
}

// AddURLS adds given objects to storage
func (st *Storage) AddURLS(ctx context.Context, urls []model.URL) ([]model.URL, error, error) {
	st.Lock()
	defer st.Unlock()

	dbObjs := schema.NewURLsFromCanonical(urls)

	for idx := range dbObjs {
		dbObjs[idx].ID = st.id

		st.urls[dbObjs[idx].ID] = dbObjs[idx]
		st.id++
	}

	objs := dbObjs.ToCanonical()

	return objs, nil, nil
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
	var urls schema.URLS
	for _, v := range st.urls {
		if v.UserID == userID {
			urls = append(urls, v)
		}
	}

	return urls.ToCanonical(), nil
}
