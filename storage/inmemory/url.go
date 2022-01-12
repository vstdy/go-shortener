package inmemory

import (
	"github.com/vstdy0/go-project/storage"
	"sync"
)

var _ storage.URLStorage = (*InMemory)(nil)

type InMemory struct {
	urls map[string]URLModel
	sync.RWMutex
}

func (inMemory *InMemory) Has(id string) bool {
	inMemory.RLock()
	defer inMemory.RUnlock()
	_, ok := inMemory.urls[id]

	return ok
}

func (inMemory *InMemory) Set(id, url string) (string, error) {
	inMemory.Lock()
	defer inMemory.Unlock()
	inMemory.urls[id] = URLModel{ID: id, URL: url}

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

func NewInMemory() *InMemory {
	var inMemory InMemory
	inMemory.urls = make(map[string]URLModel)

	return &inMemory
}
