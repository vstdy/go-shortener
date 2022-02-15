package shortener

import (
	"fmt"
	"time"

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
		delChan chan model.URL
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

	svc.delChan = make(chan model.URL)
	go svc.delWorker()

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

// delWorker starts url deletion worker.
func (svc *Service) delWorker() {
	buffer := make([]model.URL, 0, 10)

	flush := func() {
		if len(buffer) > 0 {
			err := svc.storage.RemoveUserURLs(buffer)
			if err != nil {
				log.Warn().Err(err).Msg("Objects deletion failed")
			}
		}
	}

	dur := 10 * time.Second
	timer := time.AfterFunc(dur, flush)

	for job := range svc.delChan {
		buffer = append(buffer, job)
		timer.Reset(dur)

		if cap(buffer) == len(buffer) {
			flush()
			buffer = buffer[:0]
		}
	}
}
