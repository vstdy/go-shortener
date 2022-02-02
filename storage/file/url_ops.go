package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
	inter "github.com/vstdy0/go-project/storage"
	"github.com/vstdy0/go-project/storage/file/schema"
)

var _ inter.URLStorage = (*Storage)(nil)

type (
	Storage struct {
		config  Config
		encoder *json.Encoder
		id      int
		urls    map[int]schema.URL
		sync.RWMutex
	}

	StorageOption func(st *Storage) error
)

func (st *Storage) Has(ctx context.Context, urlID int) (bool, error) {
	st.RLock()
	defer st.RUnlock()

	_, ok := st.urls[urlID]

	return ok, nil
}

func (st *Storage) Set(ctx context.Context, urls []model.URL) ([]model.URL, error) {
	st.Lock()
	defer st.Unlock()

	dbObjs := schema.NewURLsFromCanonical(urls)

	for idx := range dbObjs {
		dbObjs[idx].ID = st.id

		if err := st.encoder.Encode(dbObjs[idx]); err != nil {
			return nil, err
		}

		st.urls[dbObjs[idx].ID] = dbObjs[idx]
		st.id++
	}

	objs := dbObjs.ToCanonical()

	return objs, nil
}

func (st *Storage) Get(ctx context.Context, urlID int) (model.URL, error) {
	st.RLock()
	defer st.RUnlock()

	url, ok := st.urls[urlID]
	if !ok {
		return model.URL{}, fmt.Errorf("url does not exist")
	}

	return url.ToCanonical(), nil
}

func (st *Storage) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error) {
	var urls schema.URLS
	for _, v := range st.urls {
		if v.UserID == userID {
			urls = append(urls, v)
		}
	}

	return urls.ToCanonical(), nil
}

// WithConfig overrides default Storage config.
func WithConfig(config Config) StorageOption {
	return func(st *Storage) error {
		st.config = config

		return nil
	}
}

func New(opts ...StorageOption) (*Storage, error) {
	st := &Storage{
		config: NewDefaultConfig(),
	}
	for optIdx, opt := range opts {
		if err := opt(st); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	if err := st.config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	file, err := os.OpenFile(st.config.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	var urlModel schema.URL
	st.urls = make(map[int]schema.URL)
	st.encoder = json.NewEncoder(file)

	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		if err := json.Unmarshal(scanner.Bytes(), &urlModel); err != nil {
			return nil, err
		}
		st.urls[urlModel.ID] = urlModel
	}

	st.id = urlModel.ID + 1

	return st, nil
}
