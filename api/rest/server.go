package rest

import (
	"net/http"

	"github.com/vstdy/go-shortener/service/shortener/v1"
)

// NewServer returns new rest server.
func NewServer(svc *shortener.Service, config Config) (*http.Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	router, err := NewRouter(svc, config)
	if err != nil {
		return nil, err
	}

	return &http.Server{Addr: config.ServerAddress, Handler: router}, nil
}
