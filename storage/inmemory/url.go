package inmemory

import (
	"github.com/vstdy0/go-project/storage"
)

var _ storage.URLStorage = (*InMemory)(nil)

type InMemory struct {
	URLS map[string]URLModel
}

func (im *InMemory) Has(id string) bool {
	_, ok := im.URLS[id]
	return ok
}

func (im *InMemory) Set(id, url string) {
	im.URLS[id] = URLModel{ID: id, URL: url}
}

func (im *InMemory) Get(id string) string {
	url, ok := im.URLS[id]
	if !ok {
		return ""
	}
	return url.URL
}
