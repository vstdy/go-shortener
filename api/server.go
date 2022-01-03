package api

import (
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/service/shortener/v1"
	"net/http"
)

func NewServer(svc *shortener.Service, cfg config.Config) *http.Server {
	router := Router(svc, cfg)
	return &http.Server{Addr: cfg.ServerAddress, Handler: router}
}
