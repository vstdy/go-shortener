package inmemory

import (
	"github.com/vstdy0/go-project/storage"
)

var _ storage.URLStorage = (*URLModel)(nil)
var urls = make(map[string]URLModel)

type URLModel struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func (um URLModel) Has(id string) bool {
	_, ok := urls[id]
	return ok
}

func (um URLModel) Set(id, url string) {
	urls[id] = URLModel{ID: id, URL: url}
}

func (um URLModel) Get(id string) string {
	url, ok := urls[id]
	if !ok {
		return ""
	}
	return url.URL
}
