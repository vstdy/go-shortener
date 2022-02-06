package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/vstdy0/go-project/pkg"
	inter "github.com/vstdy0/go-project/storage"
	"github.com/vstdy0/go-project/storage/file/schema"
)

var _ inter.URLStorage = (*Storage)(nil)

type (
	Storage struct {
		config  Config
		file    *os.File
		encoder *json.Encoder
		id      int
		urls    map[int]schema.URL
		sync.RWMutex
	}

	StorageOption func(st *Storage) error
)

// WithConfig overrides default Storage config.
func WithConfig(config Config) StorageOption {
	return func(st *Storage) error {
		st.config = config

		return nil
	}
}

// New creates a new file storage with custom options.
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

	st.urls = make(map[int]schema.URL)
	st.encoder = json.NewEncoder(file)
	st.file = file

	var urlModel schema.URL
	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		if err := json.Unmarshal(scanner.Bytes(), &urlModel); err != nil {
			return nil, err
		}
		st.urls[urlModel.ID] = urlModel
	}

	st.id = urlModel.ID + 1

	return st, nil
}

// Close closes file.
func (st *Storage) Close() error {
	if st.encoder == nil {
		return nil
	}

	return st.file.Close()
}

func (st *Storage) GetPing() error {
	return pkg.ErrNoDBConnection
}
