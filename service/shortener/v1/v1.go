package shortener

import (
	"fmt"
	"github.com/vstdy0/go-project/service/shortener"
	"github.com/vstdy0/go-project/storage"
	"github.com/vstdy0/go-project/storage/inmemory"
	"strconv"
	"sync"
)

var _ shortener.URLService = (*Service)(nil)

type Service struct {
	id       int
	InMemory storage.URLStorage
	mu       sync.Mutex
}

type Option func(*Service) error

func (s *Service) AddURL(url string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.id++
	urlID := strconv.Itoa(s.id)
	s.InMemory.Set(urlID, url)
	return urlID
}

func (s *Service) GetURL(id string) string {
	return s.InMemory.Get(id)
}

func NewService(opts ...Option) (*Service, error) {
	svc := &Service{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	if svc.InMemory == nil {
		return nil, fmt.Errorf("inMemory: nil")
	}

	return svc, nil
}

func InMemoryStorage() Option {
	return func(srv *Service) error {
		im := inmemory.InMemory{}
		im.URLS = make(map[string]inmemory.URLModel)
		srv.InMemory = &im
		return nil
	}
}
