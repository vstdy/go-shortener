package storage

import (
	"bufio"
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
	sync.RWMutex
}

func (inFile *InFile) Has(id string) bool {
	inFile.RLock()
	defer inFile.RUnlock()
	_, ok := inFile.urls[id]

	return ok
}

func (inFile *InFile) Set(id, url string) (string, error) {
	inFile.Lock()
	defer inFile.Unlock()
	inFile.urls[id] = URLModel{ID: id, URL: url}
	if err := inFile.encoder.Encode(inFile.urls[id]); err != nil {
		return "", err
	}

	return id, nil
}

func (inFile *InFile) Get(id string) string {
	inFile.RLock()
	defer inFile.RUnlock()
	url, ok := inFile.urls[id]
	if !ok {
		return ""
	}

	return url.URL
}

func NewInFile(cfg config.Config) (*InFile, int, error) {
	var inFile InFile
	var urlModel URLModel
	var id int

	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, id, err
	}

	inFile.urls = make(map[string]URLModel)
	inFile.encoder = json.NewEncoder(file)

	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		if err := json.Unmarshal(scanner.Bytes(), &urlModel); err != nil {
			return nil, 0, err
		}
		inFile.urls[urlModel.ID] = urlModel
	}

	if id, err = strconv.Atoi(urlModel.ID); urlModel.ID != "" && err != nil {
		return nil, 0, err
	}

	return &inFile, id, nil
}
