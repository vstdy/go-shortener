package storage

import (
	"bufio"
	"encoding/json"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/storage"
	"os"
	"strconv"
	"sync"
)

var _ storage.URLStorage = (*InFile)(nil)

type InFile struct {
	encoder *json.Encoder
	urls    map[string]URL
	sync.RWMutex
}

func (inFile *InFile) Has(urlID string) bool {
	inFile.RLock()
	defer inFile.RUnlock()
	_, ok := inFile.urls[urlID]

	return ok
}

func (inFile *InFile) Set(urlID, userID, url string) (string, error) {
	inFile.Lock()
	defer inFile.Unlock()
	inFile.urls[urlID] = URL{ID: urlID, UserID: userID, URL: url}
	if err := inFile.encoder.Encode(inFile.urls[urlID]); err != nil {
		return "", err
	}

	return urlID, nil
}

func (inFile *InFile) Get(urlID string) string {
	inFile.RLock()
	defer inFile.RUnlock()
	url, ok := inFile.urls[urlID]
	if !ok {
		return ""
	}

	return url.URL
}

func (inFile *InFile) GetUserURLs(userID string) []model.URL {
	var urls URLS
	for _, v := range inFile.urls {
		if v.UserID == userID {
			urls = append(urls, v)
		}
	}

	return urls.ToCanonical()
}

func NewInFile(cfg config.Config) (*InFile, int, int, error) {
	var (
		inFile   InFile
		urlModel URL
		urlID    int
		userID   int
	)

	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, 0, 0, err
	}

	inFile.urls = make(map[string]URL)
	inFile.encoder = json.NewEncoder(file)

	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		if err := json.Unmarshal(scanner.Bytes(), &urlModel); err != nil {
			return nil, 0, 0, err
		}
		inFile.urls[urlModel.ID] = urlModel
		if entryUserID, err := strconv.Atoi(urlModel.UserID); urlModel.UserID != "" && err != nil {
			return nil, 0, 0, err
		} else if urlModel.UserID != "" && entryUserID > userID {
			userID = entryUserID
		}
	}

	if urlID, err = strconv.Atoi(urlModel.ID); urlModel.ID != "" && err != nil {
		return nil, 0, 0, err
	}

	return &inFile, urlID, userID, nil
}
