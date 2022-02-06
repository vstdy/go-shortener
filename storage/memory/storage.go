package memory

import (
	"sync"

	"github.com/vstdy0/go-project/pkg"
	inter "github.com/vstdy0/go-project/storage"
	"github.com/vstdy0/go-project/storage/memory/schema"
)

var _ inter.URLStorage = (*Storage)(nil)

type Storage struct {
	id   int
	urls map[int]schema.URL
	sync.RWMutex
}

// New creates a new memory storage.
func New() (*Storage, error) {
	var st Storage
	st.urls = make(map[int]schema.URL)
	st.id = 1

	return &st, nil
}

// Close implements the storage closer interface.
func (st *Storage) Close() error {
	return nil
}

func (st *Storage) GetPing() error {
	return pkg.ErrNoDBConnection
}
