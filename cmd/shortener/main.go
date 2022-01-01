package main

import (
	"context"
	"github.com/vstdy0/go-project/api"
	"github.com/vstdy0/go-project/service/shortener/v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	svc, err := shortener.NewService(shortener.WithInMemoryStorage())
	if err != nil {
		panic(err)
	}
	srv := api.NewServer(":8080", svc)

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
