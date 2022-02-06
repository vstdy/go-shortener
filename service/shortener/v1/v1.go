package shortener

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/pkg/logging"
	"github.com/vstdy0/go-project/service/shortener"
	inter "github.com/vstdy0/go-project/storage"
)

const (
	serviceName = "shortener"
)

var _ shortener.URLService = (*Service)(nil)

type (
	Service struct {
		storage inter.URLStorage
	}

	ServiceOption func(*Service)
)

func (svc *Service) AddURL(ctx context.Context, url *model.URL) (error, error) {
	urlsModel, pgErr, err := svc.storage.AddURLS(ctx, []model.URL{*url})
	if err != nil {
		return nil, err
	}

	url.ID = urlsModel[0].ID

	return pgErr, nil
}

func (svc *Service) AddBatchURLs(ctx context.Context, urls *[]model.URL) (error, error) {
	urlsModel, pgErr, err := svc.storage.AddURLS(ctx, *urls)
	if err != nil {
		return nil, err
	}

	*urls = urlsModel

	return pgErr, nil
}

func (svc *Service) GetURL(ctx context.Context, id int) (string, error) {
	urlModel, err := svc.storage.GetURL(ctx, id)
	if err != nil {
		return "", err
	}

	return urlModel.URL, nil
}

func (svc *Service) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error) {
	urls, err := svc.storage.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (svc *Service) GetPing() error {
	if err := svc.storage.GetPing(); err != nil {
		return err
	}

	return nil
}

// WithStorage sets Storage.
func WithStorage(st inter.URLStorage) ServiceOption {
	return func(svc *Service) {
		svc.storage = st
	}
}

// New creates a new shortener service.
func New(opts ...ServiceOption) (*Service, error) {
	svc := &Service{}
	for _, opt := range opts {
		opt(svc)
	}

	if svc.storage == nil {
		return nil, fmt.Errorf("storage: nil")
	}

	return svc, nil
}

// Close closes all service dependencies.
func (svc *Service) Close() error {
	if svc.storage == nil {
		return nil
	}

	if err := svc.storage.Close(); err != nil {
		return fmt.Errorf("closing storage: %w", err)
	}

	return nil
}

// Logger returns logger with service context.
func (svc *Service) Logger() zerolog.Logger {
	logCtx := log.With().Str(logging.ServiceKey, serviceName)

	return logCtx.Logger()
}
