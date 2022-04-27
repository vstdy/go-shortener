package shortener

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vstdy/go-shortener/model"
	"github.com/vstdy/go-shortener/pkg/logging"
	"github.com/vstdy/go-shortener/service/shortener"
	inter "github.com/vstdy/go-shortener/storage"
)

const (
	serviceName = "Shortener service"
)

var _ shortener.Service = (*Service)(nil)

type (
	// Service keeps service dependencies.
	Service struct {
		sync.RWMutex

		delChan chan model.URL
		config  Config
		storage inter.Storage
	}

	// ServiceOption defines functional argument for Service constructor.
	ServiceOption func(*Service) error
)

// WithConfig sets Config.
func WithConfig(config Config) ServiceOption {
	return func(svc *Service) error {
		svc.config = config

		return nil
	}
}

// WithStorage sets Storage.
func WithStorage(st inter.Storage) ServiceOption {
	return func(svc *Service) error {
		svc.storage = st

		return nil
	}
}

// NewService creates a new shortener service.
func NewService(opts ...ServiceOption) (*Service, error) {
	svc := &Service{
		config: NewDefaultConfig(),
	}
	for optIdx, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	if err := svc.config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	if svc.storage == nil {
		return nil, fmt.Errorf("storage: nil")
	}

	svc.delChan = make(chan model.URL)
	go svc.delWorker(svc.config)

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

// Logger returns logger with service field set.
func (svc *Service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}

// delWorker starts url deletion worker.
func (svc *Service) delWorker(config Config) {
	flush := func(delObjs []model.URL) {
		if len(delObjs) == 0 {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), config.DelReqTimeout)
		defer cancel()

		err := svc.storage.RemoveUsersURLs(ctx, delObjs)
		if err != nil {
			log.Warn().Err(err).Msg("Objects deletion failed")
		}
	}

	buffer := make([]model.URL, 0, config.DelBufCap)
	timer := time.AfterFunc(config.DelBufWipeTimeout, func() {
		svc.Lock()
		defer svc.Unlock()

		go flush(buffer)
		buffer = buffer[:0]
	})

	for obj := range svc.delChan {
		svc.Lock()

		timer.Reset(config.DelBufWipeTimeout)
		buffer = append(buffer, obj)

		if cap(buffer) == len(buffer) {
			go flush(buffer)
			buffer = buffer[:0]
		}

		svc.Unlock()
	}
}
