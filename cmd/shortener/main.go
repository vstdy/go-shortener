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
	svc, err := shortener.NewService(shortener.InMemoryStorage())
	if err != nil {
		log.Fatalf("Service init: %v", err)
	}
	router := api.Router(svc)
	srv := api.NewServer(":8080", router)

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
