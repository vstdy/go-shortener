package memory

import (
	"sync"

	"github.com/vstdy/go-shortener/pkg"
	inter "github.com/vstdy/go-shortener/storage"
	"github.com/vstdy/go-shortener/storage/memory/schema"
)

var _ inter.Storage = (*Storage)(nil)

// Storage keeps memory storage dependencies.
type Storage struct {
	sync.RWMutex

	id   int
	urls map[int]schema.URL
}

// NewStorage creates a new memory Storage.
func NewStorage() (*Storage, error) {
	var st Storage
	st.urls = make(map[int]schema.URL)
	st.id = 1

	return &st, nil
}

// Close implements the storage closer interface.
func (st *Storage) Close() error {
	return nil
}

// Ping implements the storage ping interface.
func (st *Storage) Ping() error {
	return pkg.ErrNoDBConnection
}
