package api

import (
	"net/http"
	"time"

	"github.com/vstdy0/go-shortener/service/shortener/v1"
)

// NewServer returns server.
func NewServer(svc *shortener.Service, config Config, timeout time.Duration) (*http.Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	router := NewRouter(svc, config, timeout)

	return &http.Server{Addr: config.ServerAddress, Handler: router}, nil
}
