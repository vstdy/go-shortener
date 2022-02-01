package shortener

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/service/shortener"
	inter "github.com/vstdy0/go-project/storage"
)

var _ shortener.URLService = (*Service)(nil)

type (
	Service struct {
		storage inter.URLStorage
	}

	ServiceOption func(*Service)
)

func (s *Service) AddURL(ctx context.Context, userID uuid.UUID, url string) (int, error) {
	urlModel, err := s.storage.Set(ctx, model.URL{UserID: userID, URL: url})
	if err != nil {
		return 0, err
	}

	return urlModel.ID, nil
}

func (s *Service) GetURL(ctx context.Context, id int) (string, error) {
	urlModel, err := s.storage.Get(ctx, id)
	if err != nil {
		return "", err
	}

	return urlModel.URL, nil
}

func (s *Service) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error) {
	urls, err := s.storage.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

// WithStorage sets Storage.
func WithStorage(st inter.URLStorage) ServiceOption {
	return func(svc *Service) {
		svc.storage = st
	}
}

func NewService(opts ...ServiceOption) (*Service, error) {
	svc := &Service{}
	for _, opt := range opts {
		opt(svc)
	}

	if svc.storage == nil {
		return nil, fmt.Errorf("storage: nil")
	}

	return svc, nil
}
