package memory

import (
	"sync"

	"github.com/vstdy0/go-project/pkg"
	inter "github.com/vstdy0/go-project/storage"
	"github.com/vstdy0/go-project/storage/memory/schema"
)

var _ inter.URLStorage = (*Storage)(nil)

// Storage keeps memory storage dependencies.
type Storage struct {
	sync.RWMutex

	id   int
	urls map[int]schema.URL
}

// New creates a new memory Storage.
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

// Ping implements the storage ping interface.
func (st *Storage) Ping() error {
	return pkg.ErrNoDBConnection
}
