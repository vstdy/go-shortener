package storage

import (
	"encoding/json"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/storage"
	"os"
	"strconv"
	"sync"
)

var _ storage.URLStorage = (*InFile)(nil)

type InFile struct {
	encoder *json.Encoder
	urls    map[string]URLModel
	mu      sync.Mutex
}

func (inFile *InFile) Has(id string) bool {
	_, ok := inFile.urls[id]
	return ok
}

func (inFile *InFile) Set(id, url string) (string, error) {
	inFile.mu.Lock()
	defer inFile.mu.Unlock()
	inFile.urls[id] = URLModel{ID: id, URL: url}
	err := inFile.encoder.Encode(inFile.urls[id])
	if err != nil {
		return "", err
	}
	return id, nil
}

func (inFile *InFile) Get(id string) string {
	url, ok := inFile.urls[id]
	if !ok {
		return ""
	}
	return url.URL
}

func NewInFile(cfg config.Config) (*InFile, int, error) {
	var inFile InFile
	var id int
	inFile.urls = make(map[string]URLModel)
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, id, err
	}
	inFile.encoder = json.NewEncoder(file)
	if dec := json.NewDecoder(file); dec.More() {
		urlModel := URLModel{}
		for dec.More() {
			err := dec.Decode(&urlModel)
			if err != nil {
				return nil, id, err
			}
			inFile.urls[urlModel.ID] = urlModel
		}
		id, err = strconv.Atoi(urlModel.ID)
		if err != nil {
			return nil, id, err
		}
	}
	return &inFile, id, nil
}
