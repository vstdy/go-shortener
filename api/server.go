package api

import (
	"net/http"
)

func Server() *http.Server {
	return &http.Server{Addr: ":8080", Handler: Router()}
}
