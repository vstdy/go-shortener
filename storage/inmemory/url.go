package inmemory

import (
	"github.com/vstdy0/go-project/storage"
)

var _ storage.URLStorage = (*InMemory)(nil)

type InMemory struct {
	URLS map[string]URLModel
}

func (inMemory *InMemory) Has(id string) bool {
	_, ok := inMemory.URLS[id]
	return ok
}

func (inMemory *InMemory) Set(id, url string) (string, error) {
	inMemory.URLS[id] = URLModel{ID: id, URL: url}
	return id, nil
}

func (inMemory *InMemory) Get(id string) string {
	url, ok := inMemory.URLS[id]
	if !ok {
		return ""
	}
	return url.URL
}
