package inmemory

import (
	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/storage"
	"sync"
)

var _ storage.URLStorage = (*InMemory)(nil)

type InMemory struct {
	urls map[string]URL
	sync.RWMutex
}

func (inMemory *InMemory) Has(id string) bool {
	inMemory.RLock()
	defer inMemory.RUnlock()
	_, ok := inMemory.urls[id]

	return ok
}

func (inMemory *InMemory) Set(id, userID, url string) (string, error) {
	inMemory.Lock()
	defer inMemory.Unlock()
	inMemory.urls[id] = URL{ID: id, UserID: userID, URL: url}

	return id, nil
}

func (inMemory *InMemory) Get(id string) string {
	inMemory.RLock()
	defer inMemory.RUnlock()
	url, ok := inMemory.urls[id]
	if !ok {
		return ""
	}

	return url.URL
}

func (inMemory *InMemory) GetUserURLs(userID string) []model.URL {
	var urls URLS
	for _, v := range inMemory.urls {
		if v.UserID == userID {
			urls = append(urls, v)
		}
	}

	return urls.ToCanonical()
}

func NewInMemory() *InMemory {
	var inMemory InMemory
	inMemory.urls = make(map[string]URL)

	return &inMemory
}
