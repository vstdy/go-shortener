package shortener

import (
	"fmt"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/service/shortener"
	"github.com/vstdy0/go-project/storage"
	infile "github.com/vstdy0/go-project/storage/infile"
	inmemory "github.com/vstdy0/go-project/storage/inmemory"
	"strconv"
)

var _ shortener.URLService = (*Service)(nil)

type Service struct {
	id         int
	urlStorage storage.URLStorage
}

type Option func(*Service) error

func (s *Service) AddURL(url string) (string, error) {
	id, err := s.urlStorage.Set(strconv.Itoa(s.id+1), url)
	if err != nil {
		return "", err
	}
	s.id++
	return id, nil
}

func (s *Service) GetURL(id string) string {
	return s.urlStorage.Get(id)
}

func WithInMemoryStorage() Option {
	return func(srv *Service) error {
		srv.urlStorage = inmemory.NewInMemory()
		return nil
	}
}

func WithInFileStorage(cfg config.Config) Option {
	return func(srv *Service) error {
		inFile, id, err := infile.NewInFile(cfg)
		if err != nil {
			return err
		}
		srv.id = id
		srv.urlStorage = inFile
		return nil
	}
}

func NewService(opts ...Option) (*Service, error) {
	svc := &Service{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	if svc.urlStorage == nil {
		return nil, fmt.Errorf("urlStorage: nil")
	}

	return svc, nil
}
