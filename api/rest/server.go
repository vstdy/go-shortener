package rest

import (
	"net/http"

	"github.com/vstdy0/go-shortener/service/shortener/v1"
)

// NewServer returns server.
func NewServer(svc *shortener.Service, config Config) (*http.Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	router := NewRouter(svc, config)

	return &http.Server{Addr: config.ServerAddress, Handler: router}, nil
}
