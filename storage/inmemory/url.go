package inmemory

import (
	"github.com/vstdy0/go-project/storage"
	"sync"
)

var _ storage.URLStorage = (*InMemory)(nil)

type InMemory struct {
	urls map[string]URLModel
	mu   sync.Mutex
}

func (inMemory *InMemory) Has(id string) bool {
	_, ok := inMemory.urls[id]
	return ok
}

func (inMemory *InMemory) Set(id, url string) (string, error) {
	inMemory.mu.Lock()
	defer inMemory.mu.Unlock()
	inMemory.urls[id] = URLModel{ID: id, URL: url}
	return id, nil
}

func (inMemory *InMemory) Get(id string) string {
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
