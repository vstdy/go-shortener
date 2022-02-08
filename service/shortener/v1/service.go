package shortener

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
