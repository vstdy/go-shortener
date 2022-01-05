package storage

import (
	"encoding/json"
	"github.com/vstdy0/go-project/storage"
)

var _ storage.URLStorage = (*InFile)(nil)

type InFile struct {
	Encoder *json.Encoder
	URLS    map[string]URLModel
}

func (inFile *InFile) Has(id string) bool {
	_, ok := inFile.URLS[id]
	return ok
}

func (inFile *InFile) Set(id, url string) (string, error) {
	inFile.URLS[id] = URLModel{ID: id, URL: url}
	inFileObj := InFile{URLS: map[string]URLModel{id: {ID: id, URL: url}}}
	err := inFile.Encoder.Encode(inFileObj.URLS)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (inFile *InFile) Get(id string) string {
	url, ok := inFile.URLS[id]
	if !ok {
		return ""
	}
	return url.URL
}
