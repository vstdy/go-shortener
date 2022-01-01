package api

import (
	"github.com/vstdy0/go-project/service/shortener/v1"
	"net/http"
)

func NewServer(addr string, svc *shortener.Service) *http.Server {
	router := Router(svc)
	return &http.Server{Addr: addr, Handler: router}
}
