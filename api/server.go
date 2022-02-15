package api

import (
	"net/http"

	"github.com/vstdy0/go-project/cmd/shortener/cmd/common"
	"github.com/vstdy0/go-project/service/shortener/v1"
)

func NewServer(svc *shortener.Service, config common.Config) *http.Server {
	router := Router(svc, config)

	return &http.Server{Addr: config.ServerAddress, Handler: router}
}
