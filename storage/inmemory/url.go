package inmemory

import (
	"github.com/vstdy0/go-project/model"
	inter "github.com/vstdy0/go-project/storage"
	"sync"
)

var _ inter.URLStorage = (*Storage)(nil)

type Storage struct {
	urls map[string]URL
	sync.RWMutex
}

func (storage *Storage) Has(id string) bool {
	storage.RLock()
	defer storage.RUnlock()
	_, ok := storage.urls[id]

	return ok
}

func (storage *Storage) Set(id, userID, url string) (string, error) {
	storage.Lock()
	defer storage.Unlock()
	storage.urls[id] = URL{ID: id, UserID: userID, URL: url}

	return id, nil
}

func (storage *Storage) Get(id string) string {
	storage.RLock()
	defer storage.RUnlock()
	url, ok := storage.urls[id]
	if !ok {
		return ""
	}

	return url.URL
}

func (storage *Storage) GetUserURLs(userID string) []model.URL {
	var urls URLS
	for _, v := range storage.urls {
		if v.UserID == userID {
			urls = append(urls, v)
		}
	}

	return urls.ToCanonical()
}

func NewInMemory() *Storage {
	var storage Storage
	storage.urls = make(map[string]URL)

	return &storage
}
