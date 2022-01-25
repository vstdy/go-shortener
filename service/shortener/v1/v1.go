package shortener

import (
	"fmt"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/service/shortener"
	"github.com/vstdy0/go-project/storage"
	infile "github.com/vstdy0/go-project/storage/infile"
	inmemory "github.com/vstdy0/go-project/storage/inmemory"
	"strconv"
)

var _ shortener.URLService = (*Service)(nil)

type Service struct {
	urlID      int
	userID     int
	urlStorage storage.URLStorage
}

type Option func(*Service) error

func (s *Service) AddURL(userID, url string) (string, error) {
	urlID, err := s.urlStorage.Set(strconv.Itoa(s.urlID+1), userID, url)
	if err != nil {
		return "", err
	}
	s.urlID++

	return urlID, nil
}

func (s *Service) GetURL(id string) string {
	return s.urlStorage.Get(id)
}

func (s *Service) GetUserURLs(userID string) []model.URL {
	return s.urlStorage.GetUserURLs(userID)
}

func (s *Service) GetUserID() int {
	return s.userID
}

func WithInMemoryStorage() Option {
	return func(srv *Service) error {
		srv.urlStorage = inmemory.NewInMemory()

		return nil
	}
}

func WithInFileStorage(cfg config.Config) Option {
	return func(srv *Service) error {
		inFile, urlID, userID, err := infile.NewInFile(cfg)
		if err != nil {
			return err
		}
		srv.urlID = urlID
		srv.userID = userID
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
