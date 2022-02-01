package api

import (
	"net/http"

	"github.com/vstdy0/go-project/cmd/shortener/cmd/common"
	"github.com/vstdy0/go-project/service/shortener/v1"
)

func NewServer(svc *shortener.Service, cfg common.Config) *http.Server {
	router := Router(svc, cfg)

	return &http.Server{Addr: cfg.ServerAddress, Handler: router}
}
