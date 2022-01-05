package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/vstdy0/go-project/api"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/service/shortener/v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	svc, err := shortener.NewService(shortener.WithInFileStorage(cfg))
	if err != nil {
		panic(err)
	}
	srv := api.NewServer(svc, cfg)

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown failed: %v", err)
		}
		log.Println("Server stopped")

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
