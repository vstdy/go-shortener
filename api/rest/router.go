package rest

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vstdy0/go-shortener/api/grpc-gateway"
	"github.com/vstdy0/go-shortener/service/shortener"
)

// NewRouter returns router.
func NewRouter(svc shortener.Service, config Config) (chi.Router, error) {
	h := NewHandler(svc, config)
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(
			middleware.RequestID,
			middleware.RealIP,
			middleware.Logger,
			middleware.Recoverer,
			middleware.StripSlashes,
			middleware.Timeout(config.Timeout),
			gzipDecompressRequest,
			gzipCompressResponse,
			cookieAuth(config.SecretKey),
		)

		r.Post("/", h.shortenURL)
		r.Get("/{id}", h.getOriginalURL)
		r.Get("/ping", h.ping)

		r.Route("/api", func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))

			r.Post("/shorten", h.shortenURL)
			r.Post("/shorten/batch", h.shortenURLsBatch)
			r.Get("/user/urls", h.getUsersURLs)
			r.Delete("/user/urls", h.deleteUserURLs)
		})
	})

	mux, err := grpcgateway.NewMux(config.GRPCServerAddress)
	if err != nil {
		return nil, fmt.Errorf("building gRPC-gateway mux: %s", err)
	}

	r.Mount("/gw", mux)

	return r, nil
}
