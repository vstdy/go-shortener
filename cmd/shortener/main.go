package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/vstdy0/go-project/api"
	"github.com/vstdy0/go-project/cmd/root"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/service/shortener/v1"
)

func main() {
	cmd := root.NewRootCmd()
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}

	cfg, err := config.LoadEnvs()
	if err != nil {
		panic(err)
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
