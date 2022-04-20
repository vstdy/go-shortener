package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/vstdy0/go-shortener/pkg"
	inter "github.com/vstdy0/go-shortener/storage"
	"github.com/vstdy0/go-shortener/storage/file/schema"
)

var _ inter.Storage = (*Storage)(nil)

type (
	// Storage keeps file storage dependencies.
	Storage struct {
		sync.RWMutex

		config  Config
		file    *os.File
		encoder *json.Encoder
		id      int
		urls    map[int]schema.URL
	}

	// StorageOption defines functional argument for Storage constructor.
	StorageOption func(st *Storage) error
)

// WithConfig overrides default Storage config.
func WithConfig(config Config) StorageOption {
	return func(st *Storage) error {
		st.config = config

		return nil
	}
}

// NewStorage creates a new file Storage with custom options.
func NewStorage(opts ...StorageOption) (*Storage, error) {
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

// Ping implements the storage ping interface.
func (st *Storage) Ping() error {
	return pkg.ErrNoDBConnection
}
