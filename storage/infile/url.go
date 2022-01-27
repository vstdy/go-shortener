package storage

import (
	"bufio"
	"encoding/json"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/model"
	inter "github.com/vstdy0/go-project/storage"
	"os"
	"strconv"
	"sync"
)

var _ inter.URLStorage = (*Storage)(nil)

type Storage struct {
	encoder *json.Encoder
	urls    map[string]URL
	sync.RWMutex
}

func (storage *Storage) Has(urlID string) bool {
	storage.RLock()
	defer storage.RUnlock()
	_, ok := storage.urls[urlID]

	return ok
}

func (storage *Storage) Set(urlID, userID, url string) (string, error) {
	storage.Lock()
	defer storage.Unlock()
	storage.urls[urlID] = URL{ID: urlID, UserID: userID, URL: url}
	if err := storage.encoder.Encode(storage.urls[urlID]); err != nil {
		return "", err
	}

	return urlID, nil
}

func (storage *Storage) Get(urlID string) string {
	storage.RLock()
	defer storage.RUnlock()
	url, ok := storage.urls[urlID]
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

func NewInFile(cfg config.Config) (*Storage, int, int, error) {
	var (
		storage  Storage
		urlModel URL
		urlID    int
		userID   int
	)

	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, 0, 0, err
	}

	storage.urls = make(map[string]URL)
	storage.encoder = json.NewEncoder(file)

	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		if err := json.Unmarshal(scanner.Bytes(), &urlModel); err != nil {
			return nil, 0, 0, err
		}
		storage.urls[urlModel.ID] = urlModel
		if entryUserID, err := strconv.Atoi(urlModel.UserID); urlModel.UserID != "" && err != nil {
			return nil, 0, 0, err
		} else if urlModel.UserID != "" && entryUserID > userID {
			userID = entryUserID
		}
	}

	if urlID, err = strconv.Atoi(urlModel.ID); urlModel.ID != "" && err != nil {
		return nil, 0, 0, err
	}

	return &storage, urlID, userID, nil
}
