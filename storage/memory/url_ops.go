package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
	inter "github.com/vstdy0/go-project/storage"
	"github.com/vstdy0/go-project/storage/memory/schema"
)

var _ inter.URLStorage = (*Storage)(nil)

type Storage struct {
	id   int
	urls map[int]schema.URL
	sync.RWMutex
}

func (st *Storage) Has(ctx context.Context, urlID int) (bool, error) {
	st.RLock()
	defer st.RUnlock()
	_, ok := st.urls[urlID]

	return ok, nil
}

func (st *Storage) Set(ctx context.Context, urls []model.URL) ([]model.URL, error) {
	st.Lock()
	defer st.Unlock()

	dbObjs := schema.NewURLsFromCanonical(urls)

	for idx := range dbObjs {
		dbObjs[idx].ID = st.id

		st.urls[dbObjs[idx].ID] = dbObjs[idx]
		st.id++
	}

	objs := dbObjs.ToCanonical()

	return objs, nil
}

func (st *Storage) Get(ctx context.Context, urlID int) (model.URL, error) {
	st.RLock()
	defer st.RUnlock()
	url, ok := st.urls[urlID]
	if !ok {
		return model.URL{}, fmt.Errorf("url does not exist")
	}

	return url.ToCanonical(), nil
}

func (st *Storage) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error) {
	var urls schema.URLS
	for _, v := range st.urls {
		if v.UserID == userID {
			urls = append(urls, v)
		}
	}

	return urls.ToCanonical(), nil
}

func New() (*Storage, error) {
	var st Storage
	st.urls = make(map[int]schema.URL)
	st.id = 1

	return &st, nil
}
